package revision

import (
	"com.fha.gocan/business/core/app"
	"com.fha.gocan/business/core/boundary"
	"com.fha.gocan/business/core/commit"
	"com.fha.gocan/business/core/revision"
	"com.fha.gocan/business/core/scene"
	revision2 "com.fha.gocan/business/data/store/revision"
	"com.fha.gocan/foundation/web"
	"github.com/jmoiron/sqlx"
	"net/http"
)

type Handlers struct {
	Revision   revision.Core
	Scene      scene.Core
	App        app.Core
	Boundary   boundary.Core
	Commit     commit.Core
	connection *sqlx.DB
}

func NewHandlers(connection *sqlx.DB) Handlers {
	return Handlers{
		Revision:   revision.NewCore(connection),
		Scene:      scene.NewCore(connection),
		App:        app.NewCore(connection),
		Boundary:   boundary.NewCore(connection),
		Commit:     commit.NewCore(connection),
		connection: connection,
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

	hotspots, err := h.Revision.QueryAppHotspots(a, beforeTime, afterTime)

	payload := struct {
		Name     string                       `json:"name"`
		Children []revision2.HotspotHierarchy `json:"children"`
	}{
		Name:     "root",
		Children: []revision2.HotspotHierarchy{hotspots},
	}

	return web.Respond(w, payload, 200)
}

func (h *Handlers) QuerySceneHotspots(w http.ResponseWriter, r *http.Request, params map[string]string) error {
	sceneId := params["sceneId"]
	s, err := h.Scene.QueryById(sceneId)
	if err != nil {
		return err
	}

	query := r.URL.Query()
	before := query.Get("before")
	after := query.Get("after")

	hotspots, err := h.Revision.QuerySceneHotspots(s.Name, before, after, h.connection)

	payload := struct {
		Name     string                       `json:"name"`
		Children []revision2.HotspotHierarchy `json:"children"`
	}{
		Name:     "root",
		Children: []revision2.HotspotHierarchy{hotspots},
	}

	return web.Respond(w, payload, 200)
}

func (h *Handlers) QueryRevisionsTrendsById(w http.ResponseWriter, r *http.Request, params map[string]string) error {
	trendId := params["trendId"]
	trends, err := h.Revision.RevisionTrendsById(trendId)
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

func (h *Handlers) QueryRevisionsTrends(w http.ResponseWriter, r *http.Request, params map[string]string) error {
	appId := params["appId"]
	trends, err := h.Revision.RevisionTrendsByAppId(appId)
	if err != nil {
		return err
	}

	payload := struct {
		Trends []revision2.RevisionTrends `json:"trends"`
	}{
		Trends: trends,
	}

	return web.Respond(w, payload, 200)
}
