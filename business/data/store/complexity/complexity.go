package complexity

import (
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type Store struct {
	connection *sqlx.DB
}

func NewStore(connection *sqlx.DB) Store {
	return Store{connection: connection}
}

func (s Store) Create(c Complexity) (Complexity, error) {
	tx := s.connection.MustBegin()

	const q1 = `insert into complexity_analysis(id, name, entity, app_id)
  values(:id, :name, :entity, :app_id)`

	if _, err := tx.NamedExec(q1, c); err != nil {
		tx.Rollback()
		return Complexity{}, errors.Wrap(err, "Unable to insert complexity analysis")
	}

	const q2 = `insert into complexity_analysis_entries(
complexity_analysis_id, 
date,
lines,
indentations,
mean,
stdev,
max
) values(
:complexity_analysis_id,
:date,
:lines,
:indentations,
:mean,
:stdev,
:max
)`
	for _, entry := range c.Entries {
		if _, err := tx.NamedExec(q2, entry); err != nil {
			tx.Rollback()
			return Complexity{}, errors.Wrap(err, "Unable to insert complexity analysis entry")
		}
	}

	tx.Commit()
	return c, nil
}
