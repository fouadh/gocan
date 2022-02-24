package file_metrics

import (
	"com.fha.gocan/foundation/db"
	"github.com/jmoiron/sqlx"
)

type Store struct {
	connection *sqlx.DB
}

func NewStore(connection *sqlx.DB) Store {
	return Store{connection: connection}
}

func (s Store) QueryMostRecent(appId string) ([]FileContent, error) {
	const q = `
	SELECT language,
       count(*) files,
       sum(blank) blanks,
       sum(comment) "comments",
       sum(code) code
	FROM cloc
	WHERE commit_id = (SELECT id
					   from commits
					   WHERE date = (
						   SELECT MAX(date)
						   from commits))
	AND app_id = :app_id
	GROUP BY language
	ORDER BY code DESC
;
`

	data := struct {
		AppId string `db:"app_id"`
	}{
		AppId: appId,
	}

	var result []FileContent
	err := db.NamedQuerySlice(s.connection, q, data, &result)
	return result, err
}
