package app

import (
	"com.fha.gocan/business/data/store/app"
	"com.fha.gocan/business/data/store/scene"
	context "com.fha.gocan/foundation"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"time"
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

func (c Core) QuerySummary(appId string, before time.Time, after time.Time) (app.Summary, error) {
	return c.app.QuerySummary(appId, before, after)
}

func (c Core) FindAppBySceneNameAndAppName(appName string, sceneName string) (app.App, error) {
	s, err := c.scene.QueryByName(sceneName)
	if err != nil {
		return app.App{}, fmt.Errorf("unable to retrieve scene %s", sceneName)
	}

	a, err := c.app.QueryBySceneIdAndName(s.Id, appName)
	if err != nil {
		return app.App{}, fmt.Errorf("unable to retrieve app %s linked to the scene %s", appName, sceneName)
	}

	return a, nil
}

func (c Core) QueryBySceneId(sceneId string) ([]app.App, error) {
	return c.app.QueryBySceneId(sceneId)
}

func (c Core) QueryById(appId string) (app.App, error) {
	return c.app.QueryById(appId)
}

func FindAppBySceneNameAndAppName(connection *sqlx.DB, appName string, sceneName string) (app.App, error) {
	c := NewCore(connection)

	return c.FindAppBySceneNameAndAppName(appName, sceneName)
}
