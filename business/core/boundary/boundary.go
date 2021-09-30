package boundary

import (
	"com.fha.gocan/business/data/store/app"
	"com.fha.gocan/business/data/store/boundary"
	"fmt"
	"github.com/jmoiron/sqlx"
	"regexp"
)

type Core struct {
	App      app.Store
	Boundary boundary.Store
}

func NewCore(connection *sqlx.DB) Core {
	return Core{
		App:      app.NewStore(connection),
		Boundary: boundary.NewStore(connection),
	}
}

func (c Core) Create(appId string, boundaryName string, transformations []string) (boundary.Boundary, error) {
	ts, err := c.parseTransformations(transformations)
	if err != nil {
		return boundary.Boundary{}, err
	}

	nb := boundary.NewBoundary{
		Name:            boundaryName,
		AppId:           appId,
		Transformations: ts,
	}

	return c.Boundary.Create(nb)
}

func (c Core) QueryByAppId(appId string) ([]boundary.Boundary, error) {
	return c.Boundary.QueryByAppId(appId)
}

func (c Core) QueryByBoundaryId(id string) (boundary.Boundary, error) {
	return c.Boundary.QueryById(id)
}

func (c Core) parseTransformations(transformations []string) ([]boundary.NewTransformation, error) {
	a := regexp.MustCompile(`:`)
	ts := []boundary.NewTransformation{}
	for _, t := range transformations {
		transformation, err := c.parseTransformation(a, t)
		if err != nil {
			return []boundary.NewTransformation{}, err
		}
		ts = append(ts, transformation)
	}
	return ts, nil
}

func (c Core) parseTransformation(a *regexp.Regexp, t string) (boundary.NewTransformation, error) {
	cols := a.Split(t, -1)
	if len(cols) != 2 {
		return boundary.NewTransformation{}, fmt.Errorf("Transformations must respect the following pattern: [Name]:[Path]. For example, tests:/src/test")
	}
	transformation := boundary.NewTransformation{
		Name: cols[0],
		Path: cols[1],
	}
	return transformation, nil
}

func (c Core) DeleteBoundaryByName(appId string, boundaryName string) error {
	return c.Boundary.DeleteByName(appId, boundaryName)
}
