package coupling

import (
	"com.fha.gocan/business/data/store/app"
	"com.fha.gocan/business/data/store/coupling"
	"com.fha.gocan/business/data/store/scene"
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

func (c Core) Query(appId string, before time.Time, after time.Time, minimalCoupling float64, minimalRevisionsAverage int) ([]coupling.Coupling, error) {
	return c.coupling.Query(appId, before, after, minimalCoupling, minimalRevisionsAverage)
}

func (c Core) QuerySoc(appId string, before time.Time, after time.Time) ([]coupling.Soc, error) {
	return c.coupling.QuerySoc(appId, before, after)
}
