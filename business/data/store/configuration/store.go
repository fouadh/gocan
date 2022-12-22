package configuration

import (
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
