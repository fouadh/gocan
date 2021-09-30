package scene

import (
	context "com.fha.gocan/foundation"
	"fmt"
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
		return Scene{}, errors.Wrap(err, "Invalid data")
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

func (s Store) QueryByName(name string) (Scene, error) {
	const q = `
	SELECT 
		id, name
	FROM
		scenes
	WHERE
		name = :scene_name
`
	var result Scene

	data := struct {
		SceneName string `db:"scene_name"`
	}{
		SceneName: name,
	}

	rows, err := s.connection.NamedQuery(q, data)
	if err != nil {
		return Scene{}, err
	}
	if !rows.Next() {
		return Scene{}, errors.New("not found")
	}

	if err := rows.StructScan(&result); err != nil {
		return Scene{}, err
	}

	return result, nil
}

func (s Store) QueryAll() ([]Scene, error) {
	const q = `
	SELECT 
		id, name 
	FROM 
		scenes
`

	var scenes []Scene
	if err := s.connection.Select(&scenes, q); err != nil {
		return []Scene{}, errors.Wrap(err, fmt.Sprintf("Cannot fetch scenes: %s", err.Error()))
	}

	return scenes, nil
}

func (s Store) QueryById(id string) (Scene, error) {
	const q = `
	SELECT 
		id, name
	FROM
		scenes
	WHERE
		id = :scene_id
`
	var result Scene

	data := struct {
		SceneId string `db:"scene_id"`
	}{
		SceneId: id,
	}

	rows, err := s.connection.NamedQuery(q, data)
	if err != nil {
		return Scene{}, err
	}
	if !rows.Next() {
		return Scene{}, errors.New("not found")
	}

	if err := rows.StructScan(&result); err != nil {
		return Scene{}, err
	}

	return result, nil
}

func (s Store) DeleteByName(name string) error {
	const q = `DELETE FROM scenes WHERE name=:scene_name`
	data := struct {
		SceneName string `db:"scene_name"`
	}{
		SceneName: name,
	}

	_, err := s.connection.NamedExec(q, data)
	return err
}
