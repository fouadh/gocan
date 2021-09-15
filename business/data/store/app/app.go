package app

import (
	context "com.fha.gocan/foundation"
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

func (s Store) QueryBySceneId(sceneId string) ([]App, error) {
	const q = `
	SELECT 
		id, name, scene_id
	FROM
		apps
	WHERE
		scene_id = :scene_id
`
	data := struct {
		SceneId string `db:"scene_id"`
	}{
		SceneId: sceneId,
	}

	rows, err := s.connection.NamedQuery(q, data)
	if err != nil {
		return []App{}, err
	}

	results := []App{}
	for rows.Next() {
		var item App
		if err := rows.StructScan(&item); err != nil {
			return []App{}, err
		}
		results = append(results, item)
	}

	return results, nil
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

	rows, err := s.connection.NamedQuery(q, data)
	if err != nil {
		return Summary{}, err
	}

	if !rows.Next() {
		return Summary{}, errors.New("No data found for this app.")
	}

	var result Summary
	if err := rows.StructScan(&result); err != nil {
		return Summary{}, err
	}

	return result, nil
}
