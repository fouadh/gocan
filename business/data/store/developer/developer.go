package developer

import (
	"github.com/jmoiron/sqlx"
	"time"
)

type Store struct {
	connection *sqlx.DB
}

func (s Store) QueryMainDevelopers(appId string, before time.Time, after time.Time) ([]EntityDeveloper, error) {
	const q = `
	SELECT 
		entity,
	    author,
		added,
	    totalAdded,
		ownership
	FROM 
		main_developers(:app_id, :before, :after)
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

	var results []EntityDeveloper

	rows, err := s.connection.NamedQuery(q, data)
	if err != nil {
		return []EntityDeveloper{}, err
	}

	for rows.Next() {
		var item EntityDeveloper
		if err := rows.StructScan(&item); err != nil {
			return []EntityDeveloper{}, err
		}
		results = append(results, item)
	}

	return results, nil
}

func (s Store) QueryEntityEfforts(appId string, before time.Time, after time.Time) ([]EntityEffort, error) {
	const q = `
	SELECT 
		entity,
	    author,
		authorRevisions,
	    totalRevisions
	FROM 
		entity_efforts(:app_id, :before, :after)
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

	var results []EntityEffort

	rows, err := s.connection.NamedQuery(q, data)
	if err != nil {
		return []EntityEffort{}, err
	}

	for rows.Next() {
		var item EntityEffort
		if err := rows.StructScan(&item); err != nil {
			return []EntityEffort{}, err
		}
		results = append(results, item)
	}

	return results, nil
}

func (s Store) QueryDevelopers(appId string, before time.Time, after time.Time) ([]Developer, error) {
	const q = `
	select name, numberOfCommits
from developers(:app_id, :before, :after)
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

	var results []Developer

	rows, err := s.connection.NamedQuery(q, data)
	if err != nil {
		return []Developer{}, err
	}

	for rows.Next() {
		var item Developer
		if err := rows.StructScan(&item); err != nil {
			return []Developer{}, err
		}
		results = append(results, item)
	}

	return results, nil
}

func NewStore(connection *sqlx.DB) Store {
	return Store{connection: connection}
}
