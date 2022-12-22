package configuration

import (
	"com.fha.gocan/foundation/db"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type Store struct {
	connection *sqlx.DB
}

func NewStore(connection *sqlx.DB) Store {
	return Store{connection: connection}
}

func (s Store) CreateExclusions(appId string, exclusions []string) error {
	const q = `
	INSERT INTO exclusions(app_id, exclusion) VALUES(:app_id, :exclusion) ON CONFLICT DO NOTHING
	`

	tx := s.connection.MustBegin()
	for _, exclusion := range exclusions {
		data := struct {
			AppId     string `db:"app_id"`
			Exclusion string `db:"exclusion"`
		}{
			AppId:     appId,
			Exclusion: exclusion,
		}

		if _, err := tx.NamedExec(q, data); err != nil {
			if err := tx.Rollback(); err != nil {
				return errors.Wrap(err, "Unable to rollback")
			}
			return errors.Wrap(err, "Unable to create exclusions.")
		}
	}

	return tx.Commit()
}

func (s Store) QueryExclusions(appId string) ([]string, error) {
	const q = `
	SELECT exclusion FROM exclusions WHERE app_id = :app_id
	`

	data := struct {
		AppId string `db:"app_id"`
	}{
		AppId: appId,
	}

	type exclusion struct {
		Exclusion string `db:"exclusion"`
	}

	var rows []exclusion

	err := db.NamedQuerySlice(s.connection, q, data, &rows)

	if err != nil {
		return nil, errors.Wrap(err, "Unable to fetch the exclusions")
	}

	results := make([]string, len(rows))

	for i := range rows {
		results[i] = rows[i].Exclusion
	}
	return results, nil
}
