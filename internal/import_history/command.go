package import_history

import (
	context "com.fha.gocan/internal/platform"
	"encoding/json"
	"fmt"
	"github.com/boyter/scc/processor"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

func NewCommand(ctx *context.Context) *cobra.Command {
	var sceneName string
	var path string

	cmd := cobra.Command{
		Use:  "import-history",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ui := ctx.Ui
			datasource := ctx.DataSource
			connection, err := datasource.GetConnection()
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("The connection to the dabase could not be established: %v", err.Error()))
			}

			ui.Say("Retrieving app...")
			appName := args[0]
			appId, err := getAppId(appName, connection, sceneName)
			if err != nil {
				return err
			}

			ui.Say("Importing history...")

			commits, err := getCommits(path)
			if err != nil {
				ui.Failed("Failed to retrieve commits: " + err.Error())
				return err
			}

			ui.Say("Importing commits...")
			err = importCommits(appId, commits, connection)
			if err != nil {
				return err
			}
			ui.Ok()

			ui.Say("Retrieving application statistics...")
			stats, err := getStats(path)
			if err != nil {
				return err
			}

			ui.Say("Importing commits statistics...")
			err = importAppStats(appId, stats, connection)
			if err != nil {
				return err
			}
			ui.Ok()

			ui.Say("Analyzing code complexity...")
			err = importCloc(appId, path, connection)
			if err != nil {
				return err
			}
			ui.Ok()
			return nil
		},
	}

	cmd.Flags().StringVarP(&sceneName, "scene", "s", "", "Scene name")
	cmd.Flags().StringVarP(&path, "path", "d", ".", "App directory")
	return &cmd
}

func getAppId(appName string, connection *sqlx.DB, sceneName string) (string, error) {
	appId := ""
	if err := connection.Get(&appId, "select id from apps where name=$1 and scene_id=(select id from scenes where name=$2)", appName, sceneName); err != nil {
		return "", errors.Wrap(err, "Unable to retrieve matching app id")
	}
	if appId == "" {
		return "", errors.New("Application not found.")
	}
	return appId, nil
}

type Clocs struct {
	Files []FileInfo
}

type FileInfo struct {
	Location string
	Code     int
}

type Info struct {
	File  string
	Lines int
}

func importCloc(appId string, directory string, connection *sqlx.DB) error {
	processor.DirFilePaths = []string{
		directory,
	}
	processor.ConfigureGc()
	processor.ConfigureLazy(true)
	processor.Format = "json"
	processor.Files = true
	processor.GitIgnore = false
	processor.Complexity = false
	file, _ := ioutil.TempFile(os.TempDir(), "gocan*.txt")
	defer os.Remove(file.Name())
	processor.FileOutput = file.Name()
	processor.Exclude = []string{"node_modules", ".idea"}
	processor.Process()

	data, _ := ioutil.ReadAll(file)
	clocs := []Clocs{}
	json.Unmarshal(data, &clocs)

	info := []Info{}
	for _, c := range clocs {
		for _, fi := range c.Files {
			info = append(info, Info{
				File:  fi.Location,
				Lines: fi.Code,
			})
		}
	}

	for _, c := range info {
		connection.MustExec(
			`INSERT INTO cloc(app_id, file, lines) VALUES($1, $2, $3)`,
			appId,
			c.File,
			c.Lines,
		)

	}

	return nil
}

func importAppStats(appId string, stats []Stat, connection *sqlx.DB) error {
	txn := connection.MustBegin()

	chunkSize := 1000

	var divided [][]Stat
	for i := 0; i < len(stats); i += chunkSize {
		end := i + chunkSize

		if end > len(stats) {
			end = len(stats)
		}

		divided = append(divided, stats[i:end])
	}

	var wg sync.WaitGroup
	wg.Add(len(divided))

	for _, set := range divided {
		go func(data []Stat) {
			defer wg.Done()
			err := bulkInsertStats(&data, appId, txn)
			if err != nil {
				// todo better than that
				fmt.Printf("Bulk Insert Error: %s", err.Error())
			}
		}(set)
	}
	wg.Wait()

	return txn.Commit()
}

func bulkInsertStats(list *[]Stat, appId string, txn *sqlx.Tx) error {
	sql := getBulkInsertSQL("stats", []string{"commit_id", "file", "insertions", "deletions", "app_id"}, len(*list))
	stmt, err := txn.Prepare(sql)
	if err != nil {
		return err
	}

	var args []interface{}
	for _, s := range *list {
		args = append(args, s.CommitId)
		args = append(args, s.File)
		args = append(args, s.Insertions)
		args = append(args, s.Deletions)
		args = append(args, appId)
	}

	_, err = stmt.Exec(args...)
	if err != nil {
		return err
	}

	err = stmt.Close()
	if err != nil {
		return err
	}
	return err

}

func importCommits(appId string, commits []Commit, connection *sqlx.DB) error {
	txn := connection.MustBegin()

	chunkSize := 1000

	var divided [][]Commit
	for i := 0; i < len(commits); i += chunkSize {
		end := i + chunkSize

		if end > len(commits) {
			end = len(commits)
		}

		divided = append(divided, commits[i:end])
	}

	var wg sync.WaitGroup
	wg.Add(len(divided))

	for _, set := range divided {
		go func(data []Commit) {
			defer wg.Done()
			err := bulkInsert(&data, appId, txn)
			if err != nil {
				// todo better than that
				fmt.Printf("Bulk Insert Error: %s", err.Error())
			}
		}(set)
	}
	wg.Wait()

	return txn.Commit()

	return nil
}

func bulkInsert(list *[]Commit, appId string, txn *sqlx.Tx) error {
	sql := getBulkInsertSQL("commits", []string{"id", "author", "date", "message", "app_id"}, len(*list))
	stmt, err := txn.Prepare(sql)
	if err != nil {
		return err
	}

	var args []interface{}
	for _, c := range *list {
		args = append(args, c.Id)
		args = append(args, c.Author)
		args = append(args, c.Date.Format(time.RFC3339))
		args = append(args, c.Message)
		args = append(args, appId)
	}

	_, err = stmt.Exec(args...)
	if err != nil {
		return err
	}

	err = stmt.Close()
	if err != nil {
		return err
	}
	return err

}

func getBulkInsertSQL(table string, columns []string, rowCount int) string {
	var b strings.Builder
	var cnt int

	columnCount := len(columns)

	b.Grow(40000) // Need to calculate, I'm too lazy))

	b.WriteString("INSERT INTO " + table + "(" + strings.Join(columns, ", ") + ") VALUES ")

	for i := 0; i < rowCount; i++ {
		b.WriteString("(")
		for j := 0; j < columnCount; j++ {
			cnt++
			b.WriteString("$")
			b.WriteString(strconv.Itoa(cnt))
			if j != columnCount-1 {
				b.WriteString(", ")
			}
		}
		b.WriteString(")")
		if i != rowCount-1 {
			b.WriteString(",")
		}
	}
	return b.String()
}

type Commit struct {
	Id      string
	Author  string
	Date    time.Time
	Message string
	AppId   string
}

type GitCommit struct {
	Id      string
	Author  string
	Date    string
	Message string
}

func getCommits(path string) ([]Commit, error) {
	cmd := exec.Command("git", "log", "--date=iso", "--pretty=format:{%n  \"Id\": \"%H\",%n  \"Author\": \"%aN\",%n  \"Date\": \"%ad\",%n  \"Message\": \"%f\"%n},")
	cmd.Dir = path
	out, err := cmd.Output()
	if err != nil {
		return nil, errors.Wrap(err, "cannot run git log command")
	}

	outStr := string(out)
	if len(outStr) == 0 {
		return nil, errors.New("no output returned by the git command")
	}

	outStr = outStr[:len(outStr)-1]
	data := fmt.Sprintf("[%s]", outStr)
	data = strings.ReplaceAll(data, "\\", "\\\\")
	gitCommits := &[]GitCommit{}
	if err := json.Unmarshal([]byte(data), &gitCommits); err != nil {
		return nil, err
	}

	commits := []Commit{}
	for _, gc := range *gitCommits {
		date, _ := time.Parse("2006-01-02 15:04:05 -0700", gc.Date)
		commits = append(commits, Commit{
			Id:      gc.Id,
			Author:  string(gc.Author),
			Date:    date,
			Message: gc.Message,
		})
	}

	return commits, nil
}

type Stat struct {
	AppId      string
	CommitId   string
	Insertions int
	Deletions  int
	File       string
}

func getStats(path string) ([]Stat, error) {
	cmd := exec.Command("git", "log", "--numstat", "--format=%H")
	cmd.Dir = path
	out, err := cmd.Output()
	if err != nil {
		return nil, errors.Wrap(err, "cannot get source code stats")
	}
	outStr := string(out)
	rows := strings.Split(outStr, "\n")
	stats := []Stat{}
	var currentCommit string
	for _, row := range rows {
		if len(strings.Split(row, "\t")) > 1 {
			stats = append(stats, buildStat(currentCommit, row))
		} else if row != "" {
			currentCommit = row
		}
	}
	return stats, nil
}

func buildStat(commitId string, line string) Stat {
	cols := strings.Split(line, "\t")
	insertions, _ := strconv.Atoi(cols[0])
	deletions, _ := strconv.Atoi(cols[1])
	return Stat{
		CommitId:   commitId,
		File:       cols[2],
		Insertions: insertions,
		Deletions:  deletions,
	}
}
