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

func (s Store) BulkImport(appId string, data []Commit) error {
	txn := s.connection.MustBegin()

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
			err := bulkInsert(&data, appId, txn)
			if err != nil {
				// todo better than that
				fmt.Printf("Bulk Insert Error: %s", err.Error())
			}
		}(set)
	}
	wg.Wait()

	return txn.Commit()

	return nil
}

func bulkInsert(list *[]Commit, appId string, txn *sqlx.Tx) error {
	sql := db.GetBulkInsertSQL("commits", []string{"id", "author", "date", "message", "app_id"}, len(*list))
	stmt, err := txn.Prepare(sql)
	if err != nil {
		return err
	}

	var args []interface{}
	for _, c := range *list {
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

