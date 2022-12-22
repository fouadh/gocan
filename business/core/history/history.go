package history

import (
	"com.fha.gocan/business/data/store/app"
	"com.fha.gocan/business/data/store/cloc"
	"com.fha.gocan/business/data/store/commit"
	"com.fha.gocan/business/data/store/configuration"
	"com.fha.gocan/business/data/store/coupling"
	"com.fha.gocan/business/data/store/scene"
	"com.fha.gocan/business/data/store/stat"
	"com.fha.gocan/business/sys/git"
	"com.fha.gocan/foundation"
	"com.fha.gocan/foundation/date"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"sort"
)

type Core struct {
	scene      scene.Store
	app        app.Store
	commit     commit.Store
	stat       stat.Store
	cloc       cloc.Store
	coupling   coupling.Store
	exclusions configuration.Store
}

func NewCore(connection *sqlx.DB) Core {
	return Core{
		scene:      scene.NewStore(connection),
		app:        app.NewStore(connection),
		commit:     commit.NewStore(connection),
		stat:       stat.NewStore(connection),
		cloc:       cloc.NewStore(connection),
		coupling:   coupling.NewStore(connection),
		exclusions: configuration.NewStore(connection),
	}
}

func (c Core) Import(appId string, path string, beforeDate string, afterDate string, beforeCommit string, afterCommit string, ctx foundation.Context, intervalBetweenAnalyses int) error {
	commits, err := git.GetCommits(path, beforeDate, afterDate, beforeCommit, afterCommit, ctx)

	sort.Slice(commits, func(i, j int) bool {
		return commits[i].Date.After(commits[j].Date)
	})

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

	exclusions, err := c.exclusions.QueryExclusionsAsGlobs(appId)
	if err != nil {
		return err
	}

	stats, err := git.GetStats(path, beforeDate, afterDate, commitsMap, ctx, exclusions)
	if err != nil {
		return err
	}

	if err = c.stat.BulkImport(appId, stats, ctx); err != nil {
		return errors.Wrap(err, "Unable to save stats")
	}

	if len(commits) == 0 {
		return errors.Errorf("No commit to analyze")
	}

	if err = c.cloc.ImportCloc(appId, path, commits[0], ctx, exclusions); err != nil {
		return errors.Wrap(err, "Unable to save clocs")
	}

	if intervalBetweenAnalyses > 0 {
		for i := len(commits) - 1; i >= 0; i -= intervalBetweenAnalyses {
			ctx.Ui.Log("Analyzing commit " + commits[i].Id + " of " + date.FormatDay(commits[i].Date))
			if err = c.cloc.ImportCloc(appId, path, commits[i], ctx, nil); err != nil {
				return errors.Wrap(err, "Unable to save clocs")
			}
		}
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
