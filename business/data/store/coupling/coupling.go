package coupling

import (
	"github.com/jmoiron/sqlx"
	"time"
)

type Store struct {
	connection *sqlx.DB
}

func (s Store) ImportCoupling(appId string, couplings []Coupling) error {
	const q = `
	INSERT INTO 
		couplings(entity, coupled, degree, averageRevisions, app_id)
	VALUES(:entity, :coupled, :degree, :average_revisions, :app_id)
`

	for _, c := range couplings {
		data := struct {
			Entity           string  `db:"entity"`
			Coupled          string  `db:"coupled"`
			Degree           float64 `db:"degree"`
			AverageRevisions float64 `db:"average_revisions"`
			AppId            string  `db:"app_id"`
		}{
			Entity:           c.Entity,
			Coupled:          c.Coupled,
			Degree:           c.Degree,
			AverageRevisions: c.AverageRevisions,
			AppId:            appId,
		}
		s.connection.NamedExec(q, data)
	}

	return nil
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
