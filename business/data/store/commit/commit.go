package commit

import (
	"com.fha.gocan/foundation/db"
	"fmt"
	"github.com/jmoiron/sqlx"
	"sync"
	"time"
)

type Store struct {
	connection *sqlx.DB
}

func NewStore(connection *sqlx.DB) Store {
	return Store{connection: connection}
}

func (s Store) Query(appId string, before time.Time, after time.Time) ([]Commit, error) {
	const q = `
	SELECT 
		id, app_id, author, message, date
	FROM
		commits
	WHERE
		app_id = :app_id
		AND date between :after and :before
`

	data := struct {
		AppId  string    `db:"app_id"`
		Before time.Time `db:"before"`
		After  time.Time `db:"after"`
	}{
		AppId: appId,
		Before: before,
		After: after,
	}

	var results []Commit
	err := db.NamedQuerySlice(s.connection, q, data, &results)
	return results, err
}

func (s Store) BulkImport(appId string, data []Commit) error {

	chunkSize := 1000

	var divided [][]Commit
	for i := 0; i < len(data); i += chunkSize {
		end := i + chunkSize

		if end > len(data) {
			end = len(data)
		}

		divided = append(divided, data[i:end])
	}

	var wg sync.WaitGroup
	wg.Add(len(divided))

	for _, set := range divided {
		go func(data []Commit) {
			defer wg.Done()
			txn := s.connection.MustBegin()
			if err := bulkInsert(data, appId, txn); err != nil {
				// todo better than that
				fmt.Printf("Bulk Insert Error: %s", err.Error())
			}
			if err := txn.Commit(); err != nil {
				// todo better than that
				fmt.Printf("Commit Error: %s", err.Error())
			}
		}(set)
	}
	wg.Wait()
	return nil
}

func (s Store) QueryCommitRange(appId string) (CommitRange, error) {
	const q = `
    	select 
			min(date) min_date, 
			max(date) max_date 
		from commits 
		where app_id=:app_id
	`

	data := struct {
		AppId string `db:"app_id"`
	}{
		AppId: appId,
	}

	var result CommitRange
	err := db.NamedQueryStruct(s.connection, q, data, &result)
	return result, err
}

func bulkInsert(list []Commit, appId string, txn *sqlx.Tx) error {
	sql := db.GetBulkInsertSQL("commits", []string{"id", "author", "date", "message", "app_id"}, len(list))
	stmt, err := txn.Prepare(sql)
	if err != nil {
		return err
	}

	var args []interface{}
	for _, c := range list {
		args = append(args, c.Id)
		args = append(args, c.Author)
		args = append(args, c.Date.Format(time.RFC3339))
		args = append(args, c.Message)
		args = append(args, appId)
	}

	_, err = stmt.Exec(args...)
	if err != nil {
		return err
	}

	err = stmt.Close()
	if err != nil {
		return err
	}

	return err
}

