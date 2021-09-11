package history

import (
	"com.fha.gocan/business/data/store/app"
	"com.fha.gocan/business/data/store/cloc"
	"com.fha.gocan/business/data/store/commit"
	"com.fha.gocan/business/data/store/scene"
	"com.fha.gocan/business/data/store/stat"
	"com.fha.gocan/business/sys/git"
	context "com.fha.gocan/foundation"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
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

func (c Core) Import(ctx context.Context, appName string, sceneName string, path string) error {
	s, err := c.scene.QueryByName(sceneName)

	if err != nil {
		return fmt.Errorf("unable to retrieve scene %s", sceneName)
	}

	a, err := c.app.QueryBySceneIdAndName(s.Id, appName)

	if err != nil {
		return fmt.Errorf("unable to retrieve app %s linked to the scene %s", appName, sceneName)
	}

	commits, err := git.GetCommits(path)
	if err != nil {
		return errors.Wrap(err, "Unable to retrieve commits")
	}
	if err = c.commit.BulkImport(a.Id, commits); err != nil {
		return errors.Wrap(err, "Unable to save commits")
	}

	stats, err := git.GetStats(path)
	if err != nil {
		return err
	}

	if err = c.stat.ImportAppStats(a.Id, stats); err != nil {
		return errors.Wrap(err, "Unable to save stats")
	}

	if err = c.cloc.ImportCloc(a.Id, path); err != nil {
		return errors.Wrap(err, "Unable to save clocs")
	}

	return nil
}


