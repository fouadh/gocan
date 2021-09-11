package cloc

import (
	"encoding/json"
	"github.com/boyter/scc/processor"
	"github.com/jmoiron/sqlx"
	"io/ioutil"
	"os"
)

type Store struct {
	connection *sqlx.DB
}

func NewStore(connection *sqlx.DB) Store {
	return Store{connection: connection}
}

func (s Store) ImportCloc(appId string, directory string) error {
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
		s.connection.MustExec(
			`INSERT INTO cloc(app_id, file, lines) VALUES($1, $2, $3)`,
			appId,
			c.File,
			c.Lines,
		)

	}

	return nil
}
