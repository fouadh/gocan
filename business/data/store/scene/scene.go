package scene

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

func (s Store) Create(ctx context.Context, newScene NewScene) (Scene, error) {
	if err := ctx.Validator.Struct(newScene); err != nil {
		return Scene{}, errors.Wrap(err, "Invalid request to create a scene")
	}

	scene := Scene{
		Id:   uuid.NewUUID().String(),
		Name: newScene.Name,
	}

	if _, err := s.connection.NamedExec("insert into scenes(id, name) values(:id, :name)", scene); err != nil {
		return Scene{}, errors.Wrap(err, "Scene could not be created")
	}

	return scene, nil
}
