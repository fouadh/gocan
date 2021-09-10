package scene

import (
	context "com.fha.gocan/foundation"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/pborman/uuid"
	"github.com/pkg/errors"
	"strings"
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

func queryString(query string, args ...interface{}) string {
	query, params, err := sqlx.Named(query, args)
	if err != nil {
		return err.Error()
	}

	for _, param := range params {
		var value string
		switch v := param.(type) {
		case string:
			value = fmt.Sprintf("%q", v)
		case []byte:
			value = fmt.Sprintf("%q", string(v))
		default:
			value = fmt.Sprintf("%v", v)
		}
		query = strings.Replace(query, "?", value, 1)
	}

	query = strings.ReplaceAll(query, "\t", "")
	query = strings.ReplaceAll(query, "\n", " ")

	return strings.Trim(query, " ")
}

