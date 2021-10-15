package boundary

import (
	"encoding/json"
	"github.com/jmoiron/sqlx"
	"github.com/pborman/uuid"
	"github.com/pkg/errors"
)

type Store struct {
	connection *sqlx.DB
}

func (s Store) Create(newBoundary NewBoundary) (Boundary, error) {
	b := s.buildBoundary(newBoundary)

	const q1 = `insert into boundaries(id, name, app_id) values(:id, :name, :app_id)`
	if _, err := s.connection.NamedExec(q1, b); err != nil {
		return Boundary{}, errors.Wrap(err, "Cannot insert a new boundary")
	}

	const q2 = `insert into transformations(boundary_id, name, path) values(:boundary_id, :name, :path)`
	for _, t := range b.Transformations {
		if _, err := s.connection.NamedExec(q2, t); err != nil {
			return Boundary{}, errors.Wrap(err, "Cannot insert a new transformation")
		}
	}

	return b, nil
}

func (s Store) buildBoundary(newBoundary NewBoundary) Boundary {
	boundaryId := uuid.NewUUID().String()
	transformations := []Transformation{}
	for _, nt := range newBoundary.Transformations {
		t := Transformation{
			BoundaryId: boundaryId,
			Name:       nt.Name,
			Path:       nt.Path,
		}
		transformations = append(transformations, t)
	}

	b := Boundary{
		Id:              boundaryId,
		Name:            newBoundary.Name,
		AppId:           newBoundary.AppId,
		Transformations: transformations,
	}
	return b
}

func (s Store) QueryByAppId(appId string) ([]Boundary, error) {
	const q = `
		select row_to_json(row) as row
from (
         select *
         from boundaries_transformations
         where app_id=:app_id
     ) row;`

	data := struct {
		AppId string `db:"app_id"`
	}{
		AppId: appId,
	}

	rows, err := s.connection.NamedQuery(q, data)
	if err != nil {
		return []Boundary{}, err
	}
	defer rows.Close()

	results := []Boundary{}
	for rows.Next() {
		var row struct {
			Row string `db:"row"`
		}
		if err := rows.StructScan(&row); err != nil {
			return []Boundary{}, err
		}
		var item Boundary
		if err := json.Unmarshal([]byte(row.Row), &item); err != nil {
			return []Boundary{}, err
		}
		results = append(results, item)
	}

	return results, nil
}

func (s Store) QueryById(boundaryId string) (Boundary, error) {
	const q = `
		select row_to_json(row) as row
from (
         select *
         from boundaries_transformations
         where id=:boundary_id
     ) row;`

	data := struct {
		BoundaryId string `db:"boundary_id"`
	}{
		BoundaryId: boundaryId,
	}

	rows, err := s.connection.NamedQuery(q, data)
	if err != nil || !rows.Next() {
		return Boundary{}, err
	}
	defer rows.Close()

	var row struct {
		Row string `db:"row"`
	}
	if err := rows.StructScan(&row); err != nil {
		return Boundary{}, err
	}

	var result Boundary
	if err := json.Unmarshal([]byte(row.Row), &result); err != nil {
		return Boundary{}, err
	}

	return result, nil
}

func (s Store) QueryByAppIdAndName(appId string, boundaryName string) (Boundary, error) {
	const q = `
		select row_to_json(row) as row
from (
         select *
         from boundaries_transformations
         where app_id=:app_id and name=:boundary_name
     ) row;`

	data := struct {
		AppId        string `db:"app_id"`
		BoundaryName string `db:"boundary_name"`
	}{
		AppId:        appId,
		BoundaryName: boundaryName,
	}

	rows, err := s.connection.NamedQuery(q, data)
	if err != nil || !rows.Next() {
		return Boundary{}, err
	}
	defer rows.Close()

	var row struct {
		Row string `db:"row"`
	}
	if err := rows.StructScan(&row); err != nil {
		return Boundary{}, err
	}

	var result Boundary
	if err := json.Unmarshal([]byte(row.Row), &result); err != nil {
		return Boundary{}, err
	}

	return result, nil
}

func (s Store) DeleteByName(appId string, boundaryName string) error {
	const q = `
	DELETE FROM boundaries
	WHERE
	app_id=:app_id AND name=:boundary_name
	`

	data := struct {
		AppId        string `db:"app_id"`
		BoundaryName string `db:"boundary_name"`
	}{
		AppId:        appId,
		BoundaryName: boundaryName,
	}

	_, err := s.connection.NamedExec(q, data)
	return err
}

func NewStore(connection *sqlx.DB) Store {
	return Store{connection: connection}
}
