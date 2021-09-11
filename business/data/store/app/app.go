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

func (s Store) QueryBySceneIdAndName(sceneId string, name string) (App, error) {
	const q = `
	SELECT 
		id, name, scene_id
	FROM
		apps
	WHERE
		name = :app_name
		AND scene_id = :scene_id
`
	var result App

	data := struct {
		SceneId string `db:"scene_id"`
		AppName string `db:"app_name"`
	}{
		SceneId: sceneId,
		AppName: name,
	}

	rows, err := s.connection.NamedQuery(q, data)
	if err != nil {
		return App{}, err
	}
	if !rows.Next() {
		return App{}, errors.New("not found")
	}

	if err := rows.StructScan(&result); err != nil {
		return App{}, err
	}

	return result, nil

}
