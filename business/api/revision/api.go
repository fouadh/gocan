package revision

import (
	"com.fha.gocan/business/api"
	"com.fha.gocan/business/core/app"
	"com.fha.gocan/business/core/boundary"
	"com.fha.gocan/business/core/commit"
	"com.fha.gocan/business/core/revision"
	"com.fha.gocan/business/core/scene"
	revision2 "com.fha.gocan/business/data/store/revision"
	"com.fha.gocan/foundation/web"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"net/http"
)

type handlers struct {
	Revision   revision.Core
	Scene      scene.Core
	App        app.Core
	Boundary   boundary.Core
	Commit     commit.Core
	connection *sqlx.DB
}

func HttpMappings(connection *sqlx.DB) api.HttpMappings {
	return handlers{
		Revision:   revision.NewCore(connection),
		Scene:      scene.NewCore(connection),
		App:        app.NewCore(connection),
		Boundary:   boundary.NewCore(connection),
		Commit:     commit.NewCore(connection),
		connection: connection,
	}
}

func (h *handlers) query(w http.ResponseWriter, r *http.Request, params map[string]string) error {
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

func (h *handlers) queryHotspots(w http.ResponseWriter, r *http.Request, params map[string]string) error {
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

	boundaryName := query.Get("boundaryName")
	moduleName := query.Get("moduleName")

	var hotspots revision2.HotspotHierarchy

	if boundaryName != "" && moduleName != "" {
		b, err := h.Boundary.QueryByAppIdAndName(appId, boundaryName)
		if err != nil {
			return err
		}

		var mod = b.FindModule(moduleName)
		if mod.Name == "" {
			return errors.New("unable to retrieve module " + moduleName)
		}

		hotspots, err = h.Revision.QueryAppHotspotsByModule(a, mod, beforeTime, afterTime)
	} else if boundaryName != "" {
		b, err := h.Boundary.QueryByAppIdAndName(appId, boundaryName)
		if err != nil {
			return err
		}
		hotspots, err = h.Revision.QueryAppHotspotsForBoundary(a, b, beforeTime, afterTime)
	} else {
		hotspots, err = h.Revision.QueryAppHotspots(a, beforeTime, afterTime)
	}

	if err != nil {
		return err
	}

	payload := struct {
		Name     string                       `json:"name"`
		Children []revision2.HotspotHierarchy `json:"children"`
	}{
		Name:     "root",
		Children: []revision2.HotspotHierarchy{hotspots},
	}

	return web.Respond(w, payload, 200)
}

func (h *handlers) querySceneHotspots(w http.ResponseWriter, r *http.Request, params map[string]string) error {
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

func (h *handlers) queryRevisionsTrendsById(w http.ResponseWriter, r *http.Request, params map[string]string) error {
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

func (h *handlers) queryRevisionsTrends(w http.ResponseWriter, r *http.Request, params map[string]string) error {
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

func (h handlers) GetMappings() map[string]func(w http.ResponseWriter, r *http.Request, params map[string]string) error {
	handlers := make(map[string]func(w http.ResponseWriter, r *http.Request, params map[string]string) error)
	handlers["/scenes/:sceneId/apps/:appId/revisions"] = h.query
	handlers["/scenes/:sceneId/apps/:appId/hotspots"] = h.queryHotspots
	handlers["/scenes/:sceneId/hotspots"] = h.querySceneHotspots
	handlers["/scenes/:sceneId/apps/:appId/revisions-trends/:trendId"] = h.queryRevisionsTrendsById
	handlers["/scenes/:sceneId/apps/:appId/revisions-trends"] = h.queryRevisionsTrends
	return handlers
}
