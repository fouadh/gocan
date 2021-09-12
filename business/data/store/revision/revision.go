package revision

import (
	"github.com/jmoiron/sqlx"
)

type Store struct {
	connection *sqlx.DB
}

func NewStore(connection *sqlx.DB) Store {
	return Store{connection: connection}
}

func (s Store) QueryByAppId(appId string) ([]Revision, error) {
	const q = `
	SELECT 
		entity,
		numberOfRevisions,
		numberOfAuthors,
		normalizedNumberOfRevisions,
		code
	FROM
		revisions
	WHERE
		app_id = :app_id
`
	// todo add dates
	results := []Revision{}

	data := struct {
		AppId string `db:"app_id"`
	}{
		AppId: appId,
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
