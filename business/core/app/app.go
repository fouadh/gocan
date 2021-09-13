package app

import (
	"com.fha.gocan/business/data/store/app"
	"com.fha.gocan/business/data/store/scene"
	context "com.fha.gocan/foundation"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type Core struct {
	scene scene.Store
	app   app.Store
}

func NewCore(connection *sqlx.DB) Core {
	return Core{
		scene: scene.NewStore(connection),
		app:   app.NewStore(connection),
	}
}

func (c Core) Create(ctx context.Context, appName string, sceneName string) (app.App, error) {
	s, err := c.scene.QueryByName(sceneName)
	if err != nil {
		return app.App{}, errors.Wrap(err, "Scene not found")
	}

	newApp := app.NewApp{Name: appName, SceneId: s.Id}
	a, err := c.app.Create(ctx, newApp)

	if err != nil {
		return app.App{}, errors.Wrap(err, "create")
	}

	return a, nil
}

func (c Core) QueryBySceneName(sceneName string) ([]app.App, error) {
	s, err := c.scene.QueryByName(sceneName)
	if err != nil {
		return []app.App{}, errors.Wrap(err, "Scene not found")
	}

	return c.app.QueryBySceneId(s.Id)
}
