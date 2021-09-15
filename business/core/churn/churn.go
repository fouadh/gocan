package churn

import (
	"com.fha.gocan/business/data/store/app"
	"com.fha.gocan/business/data/store/churn"
	"com.fha.gocan/business/data/store/scene"
	"github.com/jmoiron/sqlx"
	"time"
)

type Core struct {
	scene scene.Store
	app   app.Store
	churn churn.Store
}

func (c Core) QueryCodeChurn(appId string, before time.Time, after time.Time) ([]churn.CodeChurn, error) {
	return c.churn.QueryCodeChurn(appId, before, after)
}

func NewCore(connection *sqlx.DB) Core {
	return Core{
		scene: scene.NewStore(connection),
		app:   app.NewStore(connection),
		churn: churn.NewStore(connection),
	}
}
