package developer

import (
	"com.fha.gocan/foundation/db"
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
	err := db.NamedQuerySlice(s.connection, q, data, &results)
	return results, err
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
	err := db.NamedQuerySlice(s.connection, q, data, &results)
	return results, err
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
	err := db.NamedQuerySlice(s.connection, q, data, &results)
	return results, err
}

func NewStore(connection *sqlx.DB) Store {
	return Store{connection: connection}
}
