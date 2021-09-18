package modus_operandi

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

type Store struct {
	connection *sqlx.DB
}

func (s Store) Query(appId string, before time.Time, after time.Time) ([]WordCount, error) {
	const q =  `select word, nentry as count 
		from ts_stat('select to_tsvector(''%s'', message) from commits where app_id=''%s''')
		ORDER BY nentry DESC, ndoc DESC, word;`

	var results []WordCount
	err := s.connection.Select(&results, fmt.Sprintf(q, "english", appId))
	if err != nil {
		return []WordCount{}, err
	}

	return results, nil
}

func NewStore(connection *sqlx.DB) Store {
	return Store{connection: connection}
}
