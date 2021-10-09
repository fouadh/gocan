package active_set

import (
	"com.fha.gocan/foundation/db"
	"github.com/jmoiron/sqlx"
	"time"
)

type Store struct {
	connection *sqlx.DB
}

func (s Store) QueryOpenedEntities(appId string, before time.Time, after time.Time) ([]ActiveSetStats, error) {
	const q = `with data2 as (
    with data as (
        select file,
               DATE_TRUNC('month', min(date)) as month
        from stats
                 inner join commits c on c.id = stats.commit_id
                 and c.date between :after and :before
        where stats.app_id = :app_id
          and file not like '%=>%'
        group by file
    )
    select month, count(file) as count
    from data
    group by month
)
select month as date,
       sum(count) over (order by month asc rows between unbounded preceding and current row) as count
from data2`

	data := struct {
		AppId  string    `db:"app_id"`
		Before time.Time `db:"before"`
		After  time.Time `db:"after"`
	}{
		AppId:  appId,
		Before: before,
		After:  after,
	}

	var results []ActiveSetStats
	err := db.NamedQuerySlice(s.connection, q, data, &results)
	return results, err
}

func (s Store) QueryClosedEntities(appId string, before time.Time, after time.Time) ([]ActiveSetStats, error) {
	const q = `with data2 as (
    with data as (
        select file,
               DATE_TRUNC('month', max(date)) as month
        from stats
                 inner join commits c on c.id = stats.commit_id
                 and c.date between :after and :before
        where stats.app_id = :app_id
          and file not like '%=>%'
        group by file
    )
    select month, count(file) as count
    from data
    group by month
)
select month as date,
       sum(count) over (order by month asc rows between unbounded preceding and current row) as count
from data2`

	data := struct {
		AppId  string    `db:"app_id"`
		Before time.Time `db:"before"`
		After  time.Time `db:"after"`
	}{
		AppId:  appId,
		Before: before,
		After:  after,
	}

	var results []ActiveSetStats
	err := db.NamedQuerySlice(s.connection, q, data, &results)
	return results, err
}

func NewStore(connection *sqlx.DB) Store {
	return Store{connection: connection}
}


