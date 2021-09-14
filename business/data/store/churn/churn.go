package churn

import (
	"github.com/jmoiron/sqlx"
	"time"
)

type Store struct {
	connection *sqlx.DB
}

func (s Store) QueryCodeChurn(appId string, before time.Time, after time.Time) ([]CodeChurn, error) {
	const q = `
	SELECT 
		date,
	    added,
		deleted
	FROM 
		code_churn(:app_id, :before, :after)
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

	var results []CodeChurn

	rows, err := s.connection.NamedQuery(q, data)
	if err != nil {
		return []CodeChurn{}, err
	}

	for rows.Next() {
		var item CodeChurn
		if err := rows.StructScan(&item); err != nil {
			return []CodeChurn{}, err
		}
		results = append(results, item)
	}

	return results, nil
}

func NewStore(connection *sqlx.DB) Store {
	return Store{connection: connection}
}
