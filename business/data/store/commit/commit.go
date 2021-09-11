package commit

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Store struct {
	connection *sqlx.DB
}

func NewStore(connection *sqlx.DB) Store {
	return Store{connection: connection}
}

func (s Store) BulkImport(appId string, commits []Commit) error {
	txn := s.connection.MustBegin()

	chunkSize := 1000

	var divided [][]Commit
	for i := 0; i < len(commits); i += chunkSize {
		end := i + chunkSize

		if end > len(commits) {
			end = len(commits)
		}

		divided = append(divided, commits[i:end])
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
	sql := getBulkInsertSQL("commits", []string{"id", "author", "date", "message", "app_id"}, len(*list))
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

func getBulkInsertSQL(table string, columns []string, rowCount int) string {
	var b strings.Builder
	var cnt int

	columnCount := len(columns)

	b.Grow(40000) // Need to calculate, I'm too lazy))

	b.WriteString("INSERT INTO " + table + "(" + strings.Join(columns, ", ") + ") VALUES ")

	for i := 0; i < rowCount; i++ {
		b.WriteString("(")
		for j := 0; j < columnCount; j++ {
			cnt++
			b.WriteString("$")
			b.WriteString(strconv.Itoa(cnt))
			if j != columnCount-1 {
				b.WriteString(", ")
			}
		}
		b.WriteString(")")
		if i != rowCount-1 {
			b.WriteString(",")
		}
	}
	return b.String()
}
