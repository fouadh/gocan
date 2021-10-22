package stat

import (
	"com.fha.gocan/business/data/store/boundary"
	"com.fha.gocan/foundation"
	"com.fha.gocan/foundation/date"
	"com.fha.gocan/foundation/db"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
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

func (s Store) Query(appId string, before time.Time, after time.Time, period int, b boundary.Boundary) ([]StatInfo, error) {
	const q = `
	SELECT 
		date, commit_id, file
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
	 	return []StatInfo{}, errors.Wrap(err, "Unable to execute query")
	 }
	 defer rows.Close()

	 results := []StatInfo{}

	for rows.Next() {
		var item StatInfo
		if err := rows.StructScan(&item); err != nil {
			return []StatInfo{}, err
		}
		results = append(results, item)
	}

	if (period > 0) {
		results = aggregateCommitsPerPeriod(results, period)
	}

	if (b.Id != "") {
		results = aggregateCommitsPerBoundary(b, results)
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

func aggregateCommitsPerPeriod(stats []StatInfo, period int) []StatInfo {
	for i, s := range stats {
		s.CommitId = date.FormatDay(s.Date)
		stats[i] = s
	}

	return stats
}

func aggregateCommitsPerBoundary(b boundary.Boundary, stats []StatInfo) []StatInfo {
	statsByBoundaryMap := make(map[string](map[string]StatInfo))
	transformations := b.Transformations

	for _, s := range stats {
		if _, ok := statsByBoundaryMap[s.CommitId]; !ok {
			statsByBoundaryMap[s.CommitId] = make(map[string]StatInfo)
		}
		commits := statsByBoundaryMap[s.CommitId]

		var transformation string
		for _, t := range transformations {
			if strings.HasPrefix(s.File, t.Path) {
				transformation = t.Name
				break
			}
		}

		if transformation != "" {
			commits[transformation] = StatInfo{
				Date:     s.Date,
				CommitId: s.CommitId,
				File:     transformation,
			}
		}
	}

	statsByBoundary := []StatInfo{}
	for _, commits := range statsByBoundaryMap {
		for _, s := range commits {
			statsByBoundary = append(statsByBoundary, s)
		}
	}
	return statsByBoundary
}
