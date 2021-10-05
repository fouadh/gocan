package history

import (
	"com.fha.gocan/business/data/store/app"
	"com.fha.gocan/business/data/store/cloc"
	"com.fha.gocan/business/data/store/commit"
	"com.fha.gocan/business/data/store/coupling"
	"com.fha.gocan/business/data/store/scene"
	"com.fha.gocan/business/data/store/stat"
	"com.fha.gocan/business/sys/git"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"time"
)

type Core struct {
	scene    scene.Store
	app      app.Store
	commit   commit.Store
	stat     stat.Store
	cloc     cloc.Store
	coupling coupling.Store
}

func NewCore(connection *sqlx.DB) Core {
	return Core{
		scene:    scene.NewStore(connection),
		app:      app.NewStore(connection),
		commit:   commit.NewStore(connection),
		stat:     stat.NewStore(connection),
		cloc:     cloc.NewStore(connection),
		coupling: coupling.NewStore(connection),
	}
}

func (c Core) Import(appId string, path string, before time.Time, after time.Time) error {
	commits, err := git.GetCommits(path, before, after)
	if err != nil {
		return errors.Wrap(err, "Unable to retrieve commits")
	}
	if err = c.commit.BulkImport(appId, commits); err != nil {
		return errors.Wrap(err, "Unable to save commits")
	}

	commitsMap := make(map[string]commit.Commit)
	for _, ct := range commits {
		commitsMap[ct.Id] = ct
	}
	stats, err := git.GetStats(path, before, after, commitsMap)
	if err != nil {
		return err
	}

	if err = c.stat.BulkImport(appId, stats); err != nil {
		return errors.Wrap(err, "Unable to save stats")
	}

	if err = c.cloc.ImportCloc(appId, path, commits); err != nil {
		return errors.Wrap(err, "Unable to save clocs")
	}

	return nil
}

func (c Core) CheckIfCanImport(path string) error {
	ok, err := git.CheckIfAllCommited(path)
	if err != nil {
		return errors.Wrap(err, "Unable to check repo status")
	}
	if !ok {
		return errors.Errorf("The directory seems to contain files that have not been commited: please stash them or commit them before running the command.")
	}
	return nil
}



