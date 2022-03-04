package active_set

import (
	"com.fha.gocan/business/api"
	active_set "com.fha.gocan/business/core/active-set"
	"com.fha.gocan/business/core/commit"
	active_set2 "com.fha.gocan/business/data/store/active-set"
	"com.fha.gocan/foundation/web"
	"github.com/jmoiron/sqlx"
	"net/http"
)

type handlers struct {
	ActiveSet active_set.Core
	Commit    commit.Core
}

func HttpMappings(connection *sqlx.DB) api.HttpMappings {
	return handlers{
		ActiveSet: active_set.NewCore(connection),
		Commit:    commit.NewCore(connection),
	}
}

func (h *handlers) query(w http.ResponseWriter, r *http.Request, params map[string]string) error {
	appId := params["appId"]

	query := r.URL.Query()
	beforeTime, afterTime, err := h.Commit.ExtractDateRangeFromQueryParams(appId, query)
	if err != nil {
		return err
	}

	as, err := h.ActiveSet.Query(appId, beforeTime, afterTime)
	if err != nil {
		return err
	}

	payload := struct {
		ActiveSet []active_set2.ActiveSet `json:"activeSet"`
	}{
		ActiveSet: as,
	}

	return web.Respond(w, payload, 200)
}

func (h handlers) GetMappings() map[string]func(w http.ResponseWriter, r *http.Request, params map[string]string) error {
	handlers := make(map[string]func(w http.ResponseWriter, r *http.Request, params map[string]string) error)
	handlers["/scenes/:sceneId/apps/:appId/active-set"] = h.query
	return handlers
}
