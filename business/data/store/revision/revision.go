package revision

import (
	"github.com/jmoiron/sqlx"
	"time"
)

type Store struct {
	connection *sqlx.DB
}

func NewStore(connection *sqlx.DB) Store {
	return Store{connection: connection}
}

func (s Store) QueryByAppIdAndDateRange(appId string, before time.Time, after time.Time) ([]Revision, error) {
	const q = `
	SELECT 
		entity,
		numberOfRevisions,
		numberOfAuthors,
		normalizedNumberOfRevisions,
		code
	FROM
		revisions(:app_id, :before, :after)
`
	results := []Revision{}

	data := struct {
		AppId  string `db:"app_id"`
		Before time.Time `db:"before"`
		After  time.Time `db:"after"`
	}{
		AppId:  appId,
		Before: before,
		After:  after,
	}

	rows, err := s.connection.NamedQuery(q, data)
	if err != nil {
		return []Revision{}, err
	}

	for rows.Next() {
		var item Revision
		if err := rows.StructScan(&item); err != nil {
			return []Revision{}, err
		}
		results = append(results, item)
	}

	return results, nil
}
