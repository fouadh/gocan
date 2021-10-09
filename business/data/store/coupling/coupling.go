package coupling

import (
	"com.fha.gocan/foundation/db"
	"github.com/jmoiron/sqlx"
	"time"
)

type Store struct {
	connection *sqlx.DB
}

func (s Store) Query(appId string, minimalCoupling float64, minimalRevisionsAverage int) ([]Coupling, error) {
	const q = `
	SELECT 
		entity,
	    coupled,
		degree,
	    averageRevisions
	FROM 
		couplings
	WHERE
	    app_id = :app_id
	AND degree >= :minimal_coupling
	AND averageRevisions >= :minimal_revisions_average
	ORDER BY degree DESC
`

	data := struct {
		AppId                   string  `db:"app_id"`
		MinimalCoupling         float64 `db:"minimal_coupling"`
		MinimalRevisionsAverage int     `db:"minimal_revisions_average"`
	}{
		AppId:                   appId,
		MinimalCoupling:         minimalCoupling,
		MinimalRevisionsAverage: minimalRevisionsAverage,
	}

	var results []Coupling
	err := db.NamedQuerySlice(s.connection, q, data, &results)
	return results, err
}

func (s Store) QuerySoc(appId string, before time.Time, after time.Time) ([]Soc, error) {
	const q = `
	SELECT 
		entity,
	    soc
	FROM 
		soc(:app_id, :before, :after)
`

	data := struct {
		AppId  string    `db:"app_id"`
		Before time.Time `db:"before"`
		After  time.Time `db:"after"`
	}{
		AppId:  appId,
		Before: before,
		After:  after,
	}

	var results []Soc
	err := db.NamedQuerySlice(s.connection, q, data, &results)
	return results, err
}

func NewStore(connection *sqlx.DB) Store {
	return Store{connection: connection}
}
