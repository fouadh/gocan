package stat

import (
	"com.fha.gocan/foundation"
	"com.fha.gocan/foundation/db"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"sync"
	"time"
)

type Store struct {
	connection *sqlx.DB
}

func NewStore(connection *sqlx.DB) Store {
	return Store{connection: connection}
}

func (s Store) Query(appId string, before time.Time, after time.Time) ([]Stat, error) {
	const q = `
	SELECT 
		s.app_id, commit_id, insertions, deletions, file
	FROM
		stats s
		INNER JOIN commits c ON c.id=s.commit_id AND c.date between :after and :before
	WHERE
		s.app_id = :app_id
		AND s.file not like '%%=>%%'
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

	 rows, err := s.connection.NamedQuery(q, data)
	 if err != nil {
	 	return []Stat{}, errors.Wrap(err, "Unable to execute query")
	 }
	 defer rows.Close()

	 results := []Stat{}

	for rows.Next() {
		var item Stat
		if err := rows.StructScan(&item); err != nil {
			return []Stat{}, err
		}
		results = append(results, item)
	}

	return results, nil
}

func (s Store) BulkImport(appId string, data []Stat, ctx foundation.Context) error {
	txn := s.connection.MustBegin()

	chunkSize := 1000

	var divided [][]Stat
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
		go func(data []Stat) {
			defer wg.Done()
			err := bulkInsertStats(&data, appId, txn)
			if err != nil {
				// todo better than that
				ctx.Ui.Failed("Bulk Insert Error: " + err.Error())
			}
		}(set)
	}
	wg.Wait()

	return txn.Commit()
}

func bulkInsertStats(list *[]Stat, appId string, txn *sqlx.Tx) error {
	sql := db.GetBulkInsertSQL("stats", []string{"commit_id", "file", "insertions", "deletions", "app_id"}, len(*list))
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
