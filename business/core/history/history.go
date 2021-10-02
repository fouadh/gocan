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
	scene  scene.Store
	app    app.Store
	commit commit.Store
	stat stat.Store
	cloc cloc.Store
}

func NewCore(connection *sqlx.DB) Core {
	return Core{
		scene:  scene.NewStore(connection),
		app:    app.NewStore(connection),
		commit: commit.NewStore(connection),
		stat:   stat.NewStore(connection),
		cloc:   cloc.NewStore(connection),
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

func BuildCoupling(stats []stat.Stat) []coupling.Coupling {
	s1 := stats[0]
	s1Revs := 1
	s2 := stats[1]
	s2Revs := 1
	degree := 1.
	return []coupling.Coupling{
		{
			Entity:           s1.File,
			Coupled:          s2.File,
			Degree:           degree,
			AverageRevisions: float64(s1Revs + s2Revs) / 2,
		},
	}
}
