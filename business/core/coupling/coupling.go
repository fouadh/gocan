package coupling

import (
	"com.fha.gocan/business/data/store/app"
	"com.fha.gocan/business/data/store/coupling"
	"com.fha.gocan/business/data/store/scene"
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

type Core struct {
	scene    scene.Store
	app      app.Store
	coupling coupling.Store
}

func NewCore(connection *sqlx.DB) Core {
	return Core{
		scene:    scene.NewStore(connection),
		app:      app.NewStore(connection),
		coupling: coupling.NewStore(connection),
	}
}

func (c Core) Query(sceneName string, appName string, before time.Time, after time.Time, minimalCoupling float64, minimalRevisionsAverage int) ([]coupling.Coupling, error) {
	s, err := c.scene.QueryByName(sceneName)
	if err != nil {
		return []coupling.Coupling{}, fmt.Errorf("unable to retrieve scene %s", sceneName)
	}

	a, err := c.app.QueryBySceneIdAndName(s.Id, appName)
	if err != nil {
		return []coupling.Coupling{}, fmt.Errorf("unable to retrieve app %s linked to the scene %s", appName, sceneName)
	}

	return c.coupling.Query(a.Id, before, after, minimalCoupling, minimalRevisionsAverage)
}

func (c Core) QuerySoc(sceneName string, appName string, before time.Time, after time.Time) ([]coupling.Soc, error) {
	s, err := c.scene.QueryByName(sceneName)
	if err != nil {
		return []coupling.Soc{}, fmt.Errorf("unable to retrieve scene %s", sceneName)
	}

	a, err := c.app.QueryBySceneIdAndName(s.Id, appName)
	if err != nil {
		return []coupling.Soc{}, fmt.Errorf("unable to retrieve app %s linked to the scene %s", appName, sceneName)
	}

	return c.coupling.QuerySoc(a.Id, before, after)
}
