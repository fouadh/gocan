package complexity

import (
	"com.fha.gocan/foundation/db"
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

	const q1 = `insert into complexity_analyses(id, name, entity, app_id)
  values(:id, :name, :entity, :app_id)`

	if _, err := tx.NamedExec(q1, c); err != nil {
		if err := tx.Rollback(); err != nil {
			return Complexity{}, errors.Wrap(err, "Unable to rollback")
		}
		return Complexity{}, errors.Wrap(err, "Unable to insert complexity analysis")
	}

	const q2 = `insert into complexity_analyses_entries(
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
			if err := tx.Rollback(); err != nil {
				return Complexity{}, errors.Wrap(err, "Unable to rollback")
			}
			return Complexity{}, errors.Wrap(err, "Unable to insert complexity analysis entry")
		}
	}

	err := tx.Commit()
	return c, err
}

func (s Store) QueryAnalyses(appId string) ([]ComplexityAnalysisSummary, error) {
	const q = `
	SELECT 
		id, name
	FROM
		complexity_analyses
	WHERE
		app_id = :app_id
	ORDER BY 
		name ASC
`
	data := struct {
		AppId string `db:"app_id"`
	}{
		AppId: appId,
	}

	var results []ComplexityAnalysisSummary
	err := db.NamedQuerySlice(s.connection, q, data, &results)
	return results, err
}

func (s Store) QueryAnalysisEntriesById(complexityId string) ([]ComplexityEntry, error) {
	const q = `
	SELECT
		date,
		lines,
		indentations,
		mean,
		stdev,
		max
	FROM complexity_analyses_entries
	WHERE
		complexity_analysis_id = :complexity_id
    ORDER BY date ASC
	`

	data := struct {
		ComplexityId string `db:"complexity_id"`
	}{
		ComplexityId: complexityId,
	}

	var results []ComplexityEntry
	err := db.NamedQuerySlice(s.connection, q, data, &results)
	return results, err
}

func (s Store) DeleteAnalysisByName(appId string, name string) error {
	const q = `
	DELETE FROM complexity_analyses
	WHERE app_id=:app_id AND name=:complexity_name
`

	data := struct {
		AppId        string `db:"app_id"`
		ComplexityName string `db:"complexity_name"`
	}{
		AppId: appId,
		ComplexityName: name,
	}

	_, err := s.connection.NamedExec(q, data)
	return err
}
