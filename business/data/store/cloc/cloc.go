package cloc

import (
	"com.fha.gocan/business/data/store/commit"
	"encoding/json"
	"fmt"
	"github.com/boyter/scc/processor"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"os/exec"
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

	out, err := ExecuteCommand("git", []string{"rev-parse", "--abbrev-ref", "HEAD"}, directory)
	if err != nil {
		return errors.Wrap(err, "Cannot get git info")
	}

	initialBranch := strings.TrimRight(string(out), "\n")
	fmt.Println("initial branch is", initialBranch)

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

	_, err = ExecuteCommand("git", []string{"checkout", ct.Id}, directory)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("git checkout %s failed", ct.Id))
	}

	fmt.Println("git checkout", ct.Id)


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
		s.connection.MustExec(
			`INSERT INTO cloc(app_id, file, lines, commit_id) VALUES($1, $2, $3, $4)`,
			appId,
			c.File,
			c.Lines,
			ct.Id,
		)
	}

	fmt.Println("reinitialising repo to initial branch")
	_, err = ExecuteCommand("git", []string{"checkout", initialBranch}, directory)
	if err != nil {
		return errors.Wrap(err, "Cannot reinitialize initial branch")
	}

	return nil
}

func ExecuteCommand(command string, args []string, directory string) ([]byte, error) {
	cmd := exec.Command(command, args...)
	cmd.Dir = directory
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return out, nil
}
