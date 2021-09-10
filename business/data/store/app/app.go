package app

import (
	context "com.fha.gocan/foundation"
	"github.com/jmoiron/sqlx"
	"github.com/pborman/uuid"
	"github.com/pkg/errors"
)

type Store struct {
	connection *sqlx.DB
}

func NewStore(connection *sqlx.DB) Store {
	return Store{connection: connection}
}

func (s Store) Create(ctx context.Context, newApp NewApp) (App, error) {
	if err := ctx.Validator.Struct(newApp); err != nil {
		return App{}, errors.Wrap(err, "Invalid data")
	}

	a := App{
		Id: uuid.NewUUID().String(),
		Name: newApp.Name,
		SceneId: newApp.SceneId,
	}

	if _, err := s.connection.NamedExec("insert into apps(id, name, scene_id) values(:id, :name, :scene_id)", a); err != nil {
		return App{}, errors.Wrap(err, "creation")
	}

	return a, nil
}
