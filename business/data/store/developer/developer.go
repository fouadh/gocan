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

func (s Store) QueryEntityEffortsPerAuthor(appId string, before time.Time, after time.Time) ([]EntityEffortPerAuthor, error) {
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

	var results []EntityEffortPerAuthor
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

func (s Store) Rename(appId string, current string, new string) error {
	const q = `
	UPDATE commits SET author=:new_name WHERE app_id=:app_id AND author=:current_name
`

	data := struct {
		AppId       string `db:"app_id"`
		CurrentName string `db:"current_name"`
		NewName     string `db:"new_name"`
	}{
		AppId:       appId,
		CurrentName: current,
		NewName:     new,
	}

	_, err := s.connection.NamedExec(q, data)
	return err
}

func (s Store) QueryDevelopmentEffort(appId string, before time.Time, after time.Time) ([]EntityEffort, error) {
	const q = `
	select file entity, 1 - sum(contribution) effort
	from (
	       select t1.file, nc_author, nc_total, pow(cast(nc_author as float) / cast(nc_total as float), 2) contribution
	       from (
	             (
	                 select s.file, c.author, count(*) nc_author
	                 from stats s
	                          inner join commits c on c.id = s.commit_id
	                 group by s.file, c.author) t1
	                inner join
	            (select s.file, count(distinct c.author), count(*) nc_total
	             from stats s
	                      inner join commits c on c.id = s.commit_id
					where
					s.app_id=:app_id
					AND c.date between :after and :before
	             group by s.file) t2
	            on t1.file = t2.file
	                )
	   ) foo
	  group by foo.file
	  order by effort desc
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

func NewStore(connection *sqlx.DB) Store {
	return Store{connection: connection}
}
