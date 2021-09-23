package revision

import (
	"com.fha.gocan/business/core/app"
	"com.fha.gocan/business/core/boundary"
	"com.fha.gocan/business/core/commit"
	"com.fha.gocan/business/core/revision"
	revision2 "com.fha.gocan/business/data/store/revision"
	"com.fha.gocan/foundation/web"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"net/http"
)

type Handlers struct {
	Revision revision.Core
	App      app.Core
	Boundary boundary.Core
	Commit   commit.Core
}

func NewHandlers(connection *sqlx.DB) Handlers {
	return Handlers{
		Revision: revision.NewCore(connection),
		App: app.NewCore(connection),
		Boundary: boundary.NewCore(connection),
		Commit: commit.NewCore(connection),
	}
}

func (h *Handlers) Query(w http.ResponseWriter, r *http.Request, params map[string]string) error {
	appId := params["appId"]

	query := r.URL.Query()
	beforeTime, afterTime, err := h.Commit.ExtractDateRangeFromQueryParams(appId, query)
	if err != nil {
		return err
	}

	revs, err := h.Revision.Query(appId, beforeTime, afterTime)
	if err != nil {
		return err
	}

	result := struct {
		Revisions []revision2.Revision `json:"revisions"`
	}{
		Revisions: revs,
	}

	return web.Respond(w, result, 200)
}

func (h *Handlers) QueryHotspots(w http.ResponseWriter, r *http.Request, params map[string]string) error {
	appId := params["appId"]
	a, err := h.App.QueryById(appId)
	if err != nil {
		return err
	}

	query := r.URL.Query()
	beforeTime, afterTime, err := h.Commit.ExtractDateRangeFromQueryParams(appId, query)
	if err != nil {
		return err
	}

	hotspots, err := h.Revision.QueryHotspots(a, beforeTime, afterTime)

	payload := struct {
		Name     string                       `json:"name"`
		Children []revision2.HotspotHierarchy `json:"children"`
	}{
		Name:     "root",
		Children: []revision2.HotspotHierarchy{hotspots},
	}

	return web.Respond(w, payload, 200)
}

func (h *Handlers) QueryRevisionsTrends(w http.ResponseWriter, r *http.Request, params map[string]string) error {
	appId := params["appId"]
	a, err := h.App.QueryById(appId)
	if err != nil {
		return err
	}

	query := r.URL.Query()
	beforeTime, afterTime, err := h.Commit.ExtractDateRangeFromQueryParams(appId, query)
	if err != nil {
		return err
	}

	query = r.URL.Query()
	boundaryId := query.Get("boundaryId")
	b, err := h.Boundary.QueryByBoundaryId(boundaryId)
	if err != nil {
		return errors.Wrap(err, "Boundary not found")
	}

	trends, err := h.Revision.RevisionTrends(a.Id, b, beforeTime, afterTime)
	if err != nil {
		return err
	}

	payload := struct {
		Trends []revision2.RevisionTrend `json:"trends"`
	}{
		Trends: trends,
	}

	return web.Respond(w, payload, 200)
}
