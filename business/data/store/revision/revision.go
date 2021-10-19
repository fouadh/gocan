package revision

import (
	"com.fha.gocan/business/data/store/boundary"
	"com.fha.gocan/foundation/db"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
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

	data := struct {
		AppId  string    `db:"app_id"`
		Before time.Time `db:"before"`
		After  time.Time `db:"after"`
	}{
		AppId:  appId,
		Before: before,
		After:  after,
	}

	var results []Revision
	err := db.NamedQuerySlice(s.connection, q, data, &results)
	return results, err
}

func (s Store) QueryByBoundary(appId string, b boundary.Boundary, before time.Time, after time.Time) ([]Revision, error) {
	caseWhen := "case "
	for _, t := range b.Transformations {
		caseWhen += fmt.Sprintf("when file like '%s%%' then '%s'\n", t.Path, t.Name)
	}
	caseWhen += " end"

	const q = `
select category as entity,
       numberOfRevisions,
       cast(numberOfRevisions as decimal) / MAX(numberOfRevisions) OVER ()  as normalizedNumberOfRevisions,
       numberOfAuthors,
       coalesce((select sum(lines)
                 from cloc
                 where %s = category
                   and app_id = :app_id), 0) as code
from (
         select category,
                (select count(distinct c.author)
                 from stats s
                          inner join commits c
                                     on s.commit_id = c.id and c.date between :after and :before
                 where %s = category
                   and s.file not like '%%=>%%'
                   and s.app_id = :app_id) as numberOfAuthors,
                (select count(s.commit_id)
                 from stats s
                          inner join commits c
                                     on s.commit_id = c.id and c.date between :after and :before
                 where %s = category
                   and s.file not like '%%=>%%'
                   and s.app_id = :app_id) as numberOfRevisions
         from (select %s category
               from stats
                        inner join commits c on c.id = stats.commit_id
                   and c.date between :after and :before
               where stats.app_id = :app_id
                 and file not like '%%=>%%'
              ) a
         where category is not null
         group by category) b
ORDER BY numberOfRevisions DESC;
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

	query := fmt.Sprintf(q, caseWhen, caseWhen, caseWhen, caseWhen)
	var results []Revision
	err := db.NamedQuerySlice(s.connection, query, data, &results)
	return results, err
}

func (s Store) CreateTrend(trend NewRevisionTrends) error {
	tx := s.connection.MustBegin()

	const q1 = `insert into revision_trends(id, name, boundary_id) values(:id, :name, :boundary_id)`
	if _, err := tx.NamedExec(q1, trend); err != nil {
		if err := tx.Rollback(); err != nil {
			return errors.Wrap(err, "Unable to rollback after trying saving trends")
		}
		return errors.Wrap(err, "Unable to insert new trends")
	}

	const q2 = `insert into revision_trend_entries(
	id,
	revision_trend_id,
	date
) values(
         :id,
         :revision_trend_id,
         :date
)`

	const q3 = `insert into revision_trend_entry_revisions(
	entry_id,
	entity,
	number_of_revisions)
values(
    :entry_id,
	:entity,
	:number_of_revisions
)`

	for _, entry := range trend.Entries {
		if _, err := tx.NamedExec(q2, entry); err != nil {
			if err := tx.Rollback(); err != nil {
				return errors.Wrap(err, "Unable to rollback")
			}
			return errors.Wrap(err, "Unable to insert trend entry")
		}

		for _, rev := range entry.Revisions {
			if _, err := tx.NamedExec(q3, rev); err != nil {
				if err := tx.Rollback(); err != nil {
					return errors.Wrap(err, "Unable to rollback")
				}
				return errors.Wrap(err, "Unable to insert trend entry revision")
			}
		}
	}

	err := tx.Commit()
	return err
}
