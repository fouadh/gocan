package developer

import (
	"com.fha.gocan/business/data/store/app"
	"com.fha.gocan/business/data/store/developer"
	"com.fha.gocan/business/data/store/scene"
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

type Core struct {
	scene    scene.Store
	app      app.Store
	developer developer.Store
}

func (c Core) QueryMainDevelopers(sceneName string, appName string, before time.Time, after time.Time) ([]developer.EntityDeveloper, error) {
	s, err := c.scene.QueryByName(sceneName)
	if err != nil {
		return []developer.EntityDeveloper{}, fmt.Errorf("unable to retrieve scene %s", sceneName)
	}

	a, err := c.app.QueryBySceneIdAndName(s.Id, appName)
	if err != nil {
		return []developer.EntityDeveloper{}, fmt.Errorf("unable to retrieve app %s linked to the scene %s", appName, sceneName)
	}

	return c.developer.QueryMainDevelopers(a.Id, before, after)
}

func NewCore(connection *sqlx.DB) Core {
	return Core{
		scene:    scene.NewStore(connection),
		app:      app.NewStore(connection),
		developer: developer.NewStore(connection),
	}
}
