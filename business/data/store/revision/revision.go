package revision

import (
	"com.fha.gocan/business/data/store/boundary"
	"com.fha.gocan/foundation/db"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"sort"
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
	for _, t := range b.Modules {
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

	const q1 = `insert into revision_trends(id, name, boundary_id, app_id) values(:id, :name, :boundary_id, :app_id)`
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

type flatRevisionTrend struct {
	Id                string `db:"id"`
	Revision_trend_id string `db:"revision_trend_id"`
	Date              string `db:"date"`
	Entity            string `db:"entity"`
	NumberOfRevisions int    `db:"number_of_revisions"`
}

func (s Store) QueryTrendsByName(name string, appId string) (RevisionTrends, error) {
	const q = `
SELECT id, name, boundary_id FROM revision_trends
WHERE name=:name AND app_id=:app_id`

	data := struct {
		Name  string `db:"name"`
		AppId string `db:"app_id"`
	}{
		Name:  name,
		AppId: appId,
	}

	var trends RevisionTrends
	if err := db.NamedQueryStruct(s.connection, q, data, &trends); err != nil {
		return RevisionTrends{}, errors.Wrap(err, "Unable to fetch revision trends by name")
	}

	entries, err := s.QueryTrends(trends.Id)
	if err != nil {
		return RevisionTrends{}, errors.Wrap(err, "Unable to fetch revision trends entries")
	}

	trends.Entries = entries
	return trends, nil
}

func (s Store) QueryTrends(id string) ([]RevisionTrend, error) {
	const q = `
SELECT id, revision_trend_id, date, entity, number_of_revisions
FROM revision_trend_entries
INNER JOIN revision_trend_entry_revisions ON revision_trend_entry_revisions.entry_id = revision_trend_entries.id
WHERE revision_trend_id = :id
`
	data := struct {
		Id string `db:"id"`
	}{
		Id: id,
	}

	var rows []flatRevisionTrend

	if err := db.NamedQuerySlice(s.connection, q, data, &rows); err != nil {
		return []RevisionTrend{}, errors.Wrap(err, "Unable to fetch revision trends")
	}

	mapResults := make(map[string](*RevisionTrend))
	for _, row := range rows {
		if trend, ok := mapResults[row.Date]; ok {
			trend.Revisions = append(trend.Revisions, TrendRevision{
				EntryId:           row.Id,
				Entity:            row.Entity,
				NumberOfRevisions: row.NumberOfRevisions,
			})
		} else {
			revisions := []TrendRevision{
				{
					EntryId:           row.Id,
					Entity:            row.Entity,
					NumberOfRevisions: row.NumberOfRevisions,
				},
			}
			mapResults[row.Date] = &RevisionTrend{
				Date:      row.Date,
				Revisions: revisions,
			}
		}
	}

	var results []RevisionTrend
	for _, rt := range mapResults {
		results = append(results, *rt)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Date < results[j].Date
	})

	return results, nil
}

func (s Store) QueryTrendsByAppId(appId string) ([]RevisionTrends, error) {
	const q = `SELECT id, name, boundary_id FROM revision_trends WHERE app_id=:app_id ORDER BY name ASC`
	data := struct {
		AppId string `db:"app_id"`
	}{
		AppId: appId,
	}

	var results []RevisionTrends
	if err := db.NamedQuerySlice(s.connection, q, data, &results); err != nil {
		return []RevisionTrends{}, errors.Wrap(err, "Unable to fetch trends")
	}

	return results, nil
}
