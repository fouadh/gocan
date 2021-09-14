package developer

import (
	"fmt"
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
		AppId                   string    `db:"app_id"`
		Before                  time.Time `db:"before"`
		After                   time.Time `db:"after"`
	}{
		AppId:                   appId,
		Before:                  before,
		After:                   after,
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

	fmt.Println("***")
	fmt.Println(results[0])
	fmt.Println("***")
	return results, nil
}

func NewStore(connection *sqlx.DB) Store {
	return Store{connection: connection}
}