package age

import (
	"com.fha.gocan/foundation/db"
	"github.com/jmoiron/sqlx"
	"time"
)

type Store struct {
	connection *sqlx.DB
}

func NewStore(connection *sqlx.DB) Store {
	return Store{connection: connection}
}

func (s Store) QueryEntityAge(appId string, initialDate string, before time.Time, after time.Time) ([]EntityAge, error) {
	const q = `
		select 
			file as name, 
			(DATE_PART('year', cast(:initial_date as date)) - DATE_PART('year', max(date))) * 12 +
			(DATE_PART('month', cast(:initial_date as date)) - DATE_PART('month', max(date))) as age
		from stats
				 inner join commits c on c.id = stats.commit_id
		where c.app_id=:app_id
		AND date between :after and :before
		AND file not like '%=>%'
		group by file
		order by max(date)
`

	data := struct {
		AppId       string    `db:"app_id"`
		InitialDate string    `db:"initial_date"`
		Before      time.Time `db:"before"`
		After       time.Time `db:"after"`
	}{
		AppId:       appId,
		InitialDate: initialDate,
		Before:      before,
		After:       after,
	}

	var results []EntityAge
	err := db.NamedQuerySlice(s.connection, q, data, &results)
	return results, err
}
