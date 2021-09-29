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

	const q1 = `insert into complexity_analyses(id, name, entity, app_id)
  values(:id, :name, :entity, :app_id)`

	if _, err := tx.NamedExec(q1, c); err != nil {
		tx.Rollback()
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
			tx.Rollback()
			return Complexity{}, errors.Wrap(err, "Unable to insert complexity analysis entry")
		}
	}

	tx.Commit()
	return c, nil
}

func (s Store) QueryAnalyses(appId string) ([]ComplexityAnalysisSummary, error) {
	const q = `
	SELECT 
		id, name
	FROM
		complexity_analyses
	WHERE
		app_id = :app_id
`
	data := struct {
		AppId string `db:"app_id"`
	}{
		AppId: appId,
	}

	rows, err := s.connection.NamedQuery(q, data)
	if err != nil {
		return []ComplexityAnalysisSummary{}, nil
	}

	results := []ComplexityAnalysisSummary{}

	for rows.Next() {
		var item ComplexityAnalysisSummary
		if err := rows.StructScan(&item); err != nil {
			return []ComplexityAnalysisSummary{}, err
		}
		results = append(results, item)
	}

	return results, nil
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
	`

	data := struct {
		ComplexityId string `db:"complexity_id"`
	}{
		ComplexityId: complexityId,
	}

	rows, err := s.connection.NamedQuery(q, data)
	if err != nil {
		return []ComplexityEntry{}, nil
	}

	results := []ComplexityEntry{}

	for rows.Next() {
		var item ComplexityEntry
		if err := rows.StructScan(&item); err != nil {
			return []ComplexityEntry{}, err
		}
		results = append(results, item)
	}

	return results, nil
}
