package complexity

import (
	"com.fha.gocan/business/api"
	"com.fha.gocan/business/core/complexity"
	complexity2 "com.fha.gocan/business/data/store/complexity"
	"com.fha.gocan/foundation/web"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"net/http"
)

type handlers struct {
	Complexity complexity.Core
}

func HttpMappings(connection *sqlx.DB) api.HttpMappings {
	return handlers{
		Complexity: complexity.NewCore(connection),
	}
}

func (h *handlers) queryAnalyses(w http.ResponseWriter, r *http.Request, params map[string]string) error {
	appId := params["appId"]

	data, err := h.Complexity.QueryAnalyses(appId)
	if err != nil {
		return errors.Wrap(err, "Unable to query analyses list")
	}

	payload := struct {
		Analyses []complexity2.ComplexityAnalysisSummary `json:"analyses"`
	}{
		Analyses: data,
	}

	return web.Respond(w, payload, 200)
}

func (h *handlers) queryAnalysisEntriesById(w http.ResponseWriter, r *http.Request, params map[string]string) error {
	complexityId := params["complexityId"]

	data, err := h.Complexity.QueryAnalysisEntriesById(complexityId)
	if err != nil {
		return errors.Wrap(err, "Unable to query analyses by name")
	}

	payload := struct {
		Entries []complexity2.ComplexityEntry `json:"entries"`
	}{
		Entries: data,
	}

	return web.Respond(w, payload, 200)
}

func (h handlers) GetMappings() map[string]func(w http.ResponseWriter, r *http.Request, params map[string]string) error {
	handlers := make(map[string]func(w http.ResponseWriter, r *http.Request, params map[string]string) error)
	handlers["/scenes/:sceneId/apps/:appId/complexity-analyses"] = h.queryAnalyses
	handlers["/scenes/:sceneId/apps/:appId/complexity-analyses/:complexityId"] = h.queryAnalysisEntriesById
	return handlers
}
