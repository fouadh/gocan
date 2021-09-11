package stat

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"strconv"
	"strings"
	"sync"
)

type Store struct {
	connection *sqlx.DB
}

func NewStore(connection *sqlx.DB) Store {
	return Store{connection: connection}
}

func (s Store) ImportAppStats(appId string, stats []Stat) error {
	txn := s.connection.MustBegin()

	chunkSize := 1000

	var divided [][]Stat
	for i := 0; i < len(stats); i += chunkSize {
		end := i + chunkSize

		if end > len(stats) {
			end = len(stats)
		}

		divided = append(divided, stats[i:end])
	}

	var wg sync.WaitGroup
	wg.Add(len(divided))

	for _, set := range divided {
		go func(data []Stat) {
			defer wg.Done()
			err := bulkInsertStats(&data, appId, txn)
			if err != nil {
				// todo better than that
				fmt.Printf("Bulk Insert Error: %s", err.Error())
			}
		}(set)
	}
	wg.Wait()

	return txn.Commit()
}

func bulkInsertStats(list *[]Stat, appId string, txn *sqlx.Tx) error {
	sql := getBulkInsertSQL("stats", []string{"commit_id", "file", "insertions", "deletions", "app_id"}, len(*list))
	stmt, err := txn.Prepare(sql)
	if err != nil {
		return err
	}

	var args []interface{}
	for _, s := range *list {
		args = append(args, s.CommitId)
		args = append(args, s.File)
		args = append(args, s.Insertions)
		args = append(args, s.Deletions)
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

