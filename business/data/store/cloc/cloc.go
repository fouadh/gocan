package cloc

import (
	"com.fha.gocan/business/data/store/commit"
	"com.fha.gocan/business/sys/git"
	"encoding/json"
	"fmt"
	"github.com/boyter/scc/processor"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"strings"
)

type Store struct {
	connection *sqlx.DB
}

func NewStore(connection *sqlx.DB) Store {
	return Store{connection: connection}
}

func (s Store) ImportCloc(appId string, directory string, commits []commit.Commit) error {
	if len(commits) == 0 {
		return fmt.Errorf("No commit provided !")
	}

	initialBranch, err := git.GetCurrentBranch(directory)
	if err != nil {
		return err
	}

	processor.DirFilePaths = []string{
		directory,
	}
	processor.ConfigureGc()
	processor.ConfigureLazy(true)
	processor.Format = "json"
	processor.Files = true
	processor.GitIgnore = false
	processor.Complexity = false

	ct := commits[0]

	if err := git.Checkout(ct.Id, directory); err != nil {
		return err
	}

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
			var path string
			if directory != "." {
				path = strings.Replace(fi.Location, directory, "", 1)
				if strings.Index(path, "/") == 0 {
					path = strings.Replace(path, "/", "", 1)
				}
			} else {
				path = fi.Location
			}
			info = append(info, Info{
				File:  path,
				Lines: fi.Code,
			})
		}
	}

	for _, c := range info {
		s.connection.MustExec(
			`INSERT INTO cloc(app_id, file, lines, commit_id) VALUES($1, $2, $3, $4)`,
			appId,
			c.File,
			c.Lines,
			ct.Id,
		)
	}

	fmt.Println("reinitialising repo to initial branch")
	if err := git.Checkout(initialBranch, directory); err != nil {
		return errors.Wrap(err, "Cannot reinitialize initial branch")
	}

	return nil
}

