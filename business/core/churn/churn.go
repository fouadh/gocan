package churn

import (
	"com.fha.gocan/business/data/store/app"
	"com.fha.gocan/business/data/store/churn"
	"com.fha.gocan/business/data/store/scene"
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

type Core struct {
	scene scene.Store
	app   app.Store
	churn churn.Store
}

func (c Core) QueryCodeChurn(sceneName string, appName string, before time.Time, after time.Time) ([]churn.CodeChurn, error) {
	s, err := c.scene.QueryByName(sceneName)
	if err != nil {
		return []churn.CodeChurn{}, fmt.Errorf("unable to retrieve scene %s", sceneName)
	}

	a, err := c.app.QueryBySceneIdAndName(s.Id, appName)
	if err != nil {
		return []churn.CodeChurn{}, fmt.Errorf("unable to retrieve app %s linked to the scene %s", appName, sceneName)
	}

	return c.churn.QueryCodeChurn(a.Id, before, after)
}

func NewCore(connection *sqlx.DB) Core {
	return Core{
		scene: scene.NewStore(connection),
		app:   app.NewStore(connection),
		churn: churn.NewStore(connection),
	}
}
