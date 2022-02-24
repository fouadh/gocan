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

func (c Core) Create(appId string, boundaryName string, modules []string) (boundary.Boundary, error) {
	ts, err := c.parseModules(modules)
	if err != nil {
		return boundary.Boundary{}, err
	}

	nb := boundary.NewBoundary{
		Name:    boundaryName,
		AppId:   appId,
		Modules: ts,
	}

	return c.Boundary.Create(nb)
}

func (c Core) QueryByAppId(appId string) ([]boundary.Boundary, error) {
	return c.Boundary.QueryByAppId(appId)
}

func (c Core) QueryByBoundaryId(id string) (boundary.Boundary, error) {
	return c.Boundary.QueryById(id)
}

func (c Core) parseModules(modules []string) ([]boundary.NewModule, error) {
	a := regexp.MustCompile(`:`)
	ts := []boundary.NewModule{}
	for _, t := range modules {
		module, err := c.parseModule(a, t)
		if err != nil {
			return []boundary.NewModule{}, err
		}
		ts = append(ts, module)
	}
	return ts, nil
}

func (c Core) parseModule(a *regexp.Regexp, t string) (boundary.NewModule, error) {
	cols := a.Split(t, -1)
	if len(cols) != 2 {
		return boundary.NewModule{}, fmt.Errorf("Modules must respect the following pattern: [Name]:[Path]. For example, tests:/src/test")
	}
	module := boundary.NewModule{
		Name: cols[0],
		Path: cols[1],
	}
	return module, nil
}

func (c Core) DeleteBoundaryByName(appId string, boundaryName string) error {
	return c.Boundary.DeleteByName(appId, boundaryName)
}

func (c Core) QueryByAppIdAndName(appId string, name string) (boundary.Boundary, error) {
	return c.Boundary.QueryByAppIdAndName(appId, name)
}
