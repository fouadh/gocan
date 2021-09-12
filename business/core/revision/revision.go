package revision

import (
	"com.fha.gocan/business/data/store/app"
	"com.fha.gocan/business/data/store/revision"
	"com.fha.gocan/business/data/store/scene"
	context "com.fha.gocan/foundation"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"time"
)

type Core struct {
	revision revision.Store
	scene scene.Store
	app   app.Store
}

func NewCore(connection *sqlx.DB) Core {
	return Core{
		revision: revision.NewStore(connection),
		scene: scene.NewStore(connection),
		app: app.NewStore(connection),
	}
}

func (c Core) GetRevisions(ctx context.Context, appName string, sceneName string, before time.Time, after time.Time) ([]revision.Revision, error) {
	s, err := c.scene.QueryByName(sceneName)
	if err != nil {
		return []revision.Revision{}, errors.Wrap(err, "Scene not found")
	}

	a, err := c.app.QueryBySceneIdAndName(s.Id, appName)
	if err != nil {
		return []revision.Revision{}, errors.Wrap(err, "App not found")
	}

	return c.revision.QueryByAppId(a.Id)
}
