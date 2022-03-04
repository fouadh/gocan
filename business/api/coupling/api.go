package coupling

import (
	"com.fha.gocan/business/api"
	"com.fha.gocan/business/core/app"
	"com.fha.gocan/business/core/commit"
	"com.fha.gocan/business/core/coupling"
	coupling2 "com.fha.gocan/business/data/store/coupling"
	"com.fha.gocan/foundation/web"
	"github.com/jmoiron/sqlx"
	"net/http"
	"strconv"
)

type handlers struct {
	Coupling coupling.Core
	App      app.Core
	Commit   commit.Core
}

func HttpMappings(connection *sqlx.DB) api.HttpMappings {
	return handlers{
		Coupling: coupling.NewCore(connection),
		App:      app.NewCore(connection),
		Commit:   commit.NewCore(connection),
	}
}

func (h *handlers) buildCouplingHierarchy(w http.ResponseWriter, r *http.Request, params map[string]string) error {
	appId := params["appId"]

	query := r.URL.Query()
	before, after, err := h.Commit.ExtractDateRangeFromQueryParams(appId, query)
	if err != nil {
		return err
	}

	a, err := h.App.QueryById(appId)
	if err != nil {
		return err
	}

	minCouplingStr := query.Get("minCoupling")
	minCoupling, err := strconv.ParseFloat(minCouplingStr, 32)
	if err != nil {
		minCoupling = .39
	}

	minRevsAvgStr := query.Get("minRevisionsAvg")
	minRevsAvg, err := strconv.Atoi(minRevsAvgStr)
	if err != nil {
		minRevsAvg = 6
	}

	boundaryName := query.Get("boundaryName")

	var c coupling2.CouplingHierarchy
	if boundaryName == "" {
		c, err = h.Coupling.BuildCouplingHierarchy(a, minCoupling, minRevsAvg, before, after)
	} else {
		c, err = h.Coupling.BuildCouplingHierarchyByBoundary(a, boundaryName, minCoupling, minRevsAvg, before, after)
	}
	if err != nil {
		return err
	}
	return web.Respond(w, c, 200)
}

func (h *handlers) buildEntityCouplingHierarchy(w http.ResponseWriter, r *http.Request, params map[string]string) error {
	appId := params["appId"]

	query := r.URL.Query()
	entity := query.Get("entity")

	before, after, err := h.Commit.ExtractDateRangeFromQueryParams(appId, query)
	if err != nil {
		return err
	}

	a, err := h.App.QueryById(appId)
	if err != nil {
		return err
	}

	minCouplingStr := query.Get("minCoupling")
	minCoupling, err := strconv.ParseFloat(minCouplingStr, 32)
	if err != nil {
		minCoupling = .39
	}

	minRevsAvgStr := query.Get("minRevisionsAvg")
	minRevsAvg, err := strconv.Atoi(minRevsAvgStr)
	if err != nil {
		minRevsAvg = 6
	}

	c, err := h.Coupling.BuildEntityCouplingHierarchy(a, entity, minCoupling, minRevsAvg, before, after)

	if err != nil {
		return err
	}

	return web.Respond(w, c, 200)
}

func (h handlers) GetMappings() map[string]func(w http.ResponseWriter, r *http.Request, params map[string]string) error {
	handlers := make(map[string]func(w http.ResponseWriter, r *http.Request, params map[string]string) error)
	handlers["/scenes/:sceneId/apps/:appId/coupling-hierarchy"] = h.buildCouplingHierarchy
	handlers["/scenes/:sceneId/apps/:appId/entity-coupling"] = h.buildEntityCouplingHierarchy
	return handlers
}
