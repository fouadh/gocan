package revision

import (
	"com.fha.gocan/business/data/store/boundary"
	"fmt"
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
		AppId  string    `db:"app_id"`
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

	results := []Revision{}

	data := struct {
		AppId    string    `db:"app_id"`
		Before   time.Time `db:"before"`
		After    time.Time `db:"after"`
	}{
		AppId:    appId,
		Before:   before,
		After:    after,
	}

	rows, err := s.connection.NamedQuery(fmt.Sprintf(q, caseWhen, caseWhen, caseWhen, caseWhen), data)
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
