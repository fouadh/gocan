package cloc

import (
	"com.fha.gocan/business/data/store/commit"
	"com.fha.gocan/business/sys/git"
	"com.fha.gocan/foundation"
	"encoding/json"
	"github.com/boyter/scc/processor"
	"github.com/gobwas/glob"
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

func include(filename string, exclusions []glob.Glob) bool {
	for _, exclusion := range exclusions {
		if exclusion.Match(filename) {
			return false
		}
	}

	return true
}

func (s Store) ImportCloc(appId string, directory string, ct commit.Commit, ctx foundation.Context, exclusions []glob.Glob) error {
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
	processor.GitIgnore = true
	processor.Complexity = false

	if err := git.Checkout(ct.Id, directory); err != nil {
		return err
	}

	file, _ := ioutil.TempFile(os.TempDir(), "gocan*.txt")
	defer os.Remove(file.Name())
	processor.FileOutput = file.Name()
	processor.Process()

	data, _ := ioutil.ReadAll(file)
	clocs := []Clocs{}
	json.Unmarshal(data, &clocs)

	for _, c := range clocs {
		for _, fi := range c.Files {
			if include(fi.Filename, exclusions) {
				var path string
				if directory != "." {
					path = strings.Replace(fi.Location, directory, "", 1)
					if strings.Index(path, "/") == 0 {
						path = strings.Replace(path, "/", "", 1)
					}
				} else {
					path = fi.Location
				}

				fi.AppId = appId
				fi.CommitId = ct.Id
				fi.Location = path

				const q = `
					insert into cloc(
									 app_id, 
									 commit_id, 
									 file, lines, 
									 language, extension, filename, code, comment, blank, complexity, is_binary
					) values(
						:app_id, 
						:commit_id, 
						:file, 
						:lines, 
						:language, 
						:extension, 
						:filename, 
						:code, 
						:comment,
						:blank, 
						:complexity, 
						:is_binary
					) ON CONFLICT DO NOTHING
`
				if _, err := s.connection.NamedExec(q, fi); err != nil {
					return errors.Wrap(err, "Unable to save cloc analysis")
				}
			} else {
				ctx.Ui.Log("Exclude " + fi.Filename + " from cloc analysis")
			}
		}
	}

	ctx.Ui.Log("reinitialising repo to initial branch")
	if err := git.Checkout(initialBranch, directory); err != nil {
		return errors.Wrap(err, "Cannot reinitialize initial branch")
	}

	return nil
}
