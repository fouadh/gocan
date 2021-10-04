package commit

import (
	"com.fha.gocan/foundation/db"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"strconv"
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
	fmt.Println("Number of chunks is " + strconv.Itoa(len(divided)))
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
	fmt.Println("done")



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

	rows, err := s.connection.NamedQuery(q, data)
	if err != nil {
		return CommitRange{}, errors.Wrap(err, "Cannot query commit range")
	}

	if !rows.Next() {
		return CommitRange{}, errors.Wrap(err, "No commit range available")
	}

	var result CommitRange
	err = rows.StructScan(&result)
	if err != nil {
		return CommitRange{}, errors.Wrap(err, "Cannot read commit range from db")
	}

	return result, nil
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

