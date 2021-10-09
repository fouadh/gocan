package app

import (
	context "com.fha.gocan/foundation"
	"com.fha.gocan/foundation/db"
	"github.com/jmoiron/sqlx"
	"github.com/pborman/uuid"
	"github.com/pkg/errors"
	"time"
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
		Id:      uuid.NewUUID().String(),
		Name:    newApp.Name,
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

	data := struct {
		SceneId string `db:"scene_id"`
		AppName string `db:"app_name"`
	}{
		SceneId: sceneId,
		AppName: name,
	}

	var result App
	err := db.NamedQueryStruct(s.connection, q, data, &result)
	return result, err
}

func (s Store) QueryBySceneId(sceneId string) ([]App, error) {
	const q = `
	SELECT 
		id, name, scene_id
	FROM
		apps
	WHERE
		scene_id = :scene_id
	ORDER BY name
`
	data := struct {
		SceneId string `db:"scene_id"`
	}{
		SceneId: sceneId,
	}

	var results []App
	err := db.NamedQuerySlice(s.connection, q, data, &results)
	return results, err
}

func (s Store) QuerySummary(appId string, before time.Time, after time.Time) (Summary, error) {
	const q = `
	SELECT 
		name,
		id                      ,
		numberOfCommits         ,
		numberOfEntities        ,
		numberOfEntitiesChanged ,
		numberOfAuthors         
	FROM
		app_summary(:app_id, :before, :after)
`
	data := struct {
		AppId string `db:"app_id"`
		Before time.Time `db:"before"`
		After  time.Time `db:"after"`
	}{
		AppId: appId,
		Before: before,
		After:  after,
	}

	var result Summary
	err := db.NamedQueryStruct(s.connection, q, data, &result)
	return result, err
}

func (s Store) QueryById(appId string) (App, error) {
	const q = `
	SELECT 
		id, name, scene_id
	FROM
		apps
	WHERE
		id = :app_id
`
	data := struct {
		AppId string `db:"app_id"`
	}{
		AppId: appId,
	}

	var result App
	err := db.NamedQueryStruct(s.connection, q, data, &result)
	return result, err
}

func (s Store) Delete(appId string) error {
	const q = `DELETE FROM apps WHERE id=:app_id`

	data := struct {
		AppId string `db:"app_id"`
	}{
		AppId: appId,
	}

	_, err := s.connection.NamedExec(q, data)
	return err
}
