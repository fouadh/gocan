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
	stat   stat.Store
	cloc   cloc.Store
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
	s2 := stats[1]

	commits := make(map[string](map[string](bool)))
	commits[stats[0].CommitId] = map[string](bool) { stats[0].File: true, stats[1].File: true }
	commits[stats[2].CommitId] = map[string](bool) { stats[2].File: true }

	count := 0
	s1Revs := 0
	s2Revs := 0

	for _, s := range stats {
		if s.File == s1.File {
			s1Revs++
		}
		if s.File == s2.File {
			s2Revs++
			files := commits[s.CommitId]
			if _, ok := files[s.File]; ok {
				count++
			}
		}

	}
	average := float64(s1Revs+s2Revs) / 2
	degree := float64(count) / average
	return []coupling.Coupling{
		{
			Entity:           s1.File,
			Coupled:          s2.File,
			Degree:           degree,
			AverageRevisions: average,
		},
	}
}
