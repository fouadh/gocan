package coupling

import (
	"github.com/jmoiron/sqlx"
	"time"
)

type Store struct {
	connection *sqlx.DB
}

func (s Store) Query(appId string, before time.Time, after time.Time, minimalCoupling float64, minimalRevisionsAverage int) ([]Coupling, error) {
	const q = `
	SELECT 
		entity,
	    coupled,
		degree,
	    averageRevisions
	FROM 
		couplings(:app_id, :before, :after, :minimal_coupling, :minimal_revisions_average)
`

	data := struct {
		AppId                   string    `db:"app_id"`
		Before                  time.Time `db:"before"`
		After                   time.Time `db:"after"`
		MinimalCoupling         float64   `db:"minimal_coupling"`
		MinimalRevisionsAverage int       `db:"minimal_revisions_average"`
	}{
		AppId:                   appId,
		Before:                  before,
		After:                   after,
		MinimalCoupling:         minimalCoupling,
		MinimalRevisionsAverage: minimalRevisionsAverage,
	}

	var results []Coupling

	rows, err := s.connection.NamedQuery(q, data)
	if err != nil {
		return []Coupling{}, err
	}

	for rows.Next() {
		var item Coupling
		if err := rows.StructScan(&item); err != nil {
			return []Coupling{}, err
		}
		results = append(results, item)
	}

	return results, nil
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

	rows, err := s.connection.NamedQuery(q, data)
	if err != nil {
		return []Soc{}, err
	}

	for rows.Next() {
		var item Soc
		if err := rows.StructScan(&item); err != nil {
			return []Soc{}, err
		}
		results = append(results, item)
	}

	return results, nil
}

func NewStore(connection *sqlx.DB) Store {
	return Store{connection: connection}
}
