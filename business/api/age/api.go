package age

import (
	"com.fha.gocan/business/api"
	"com.fha.gocan/business/core/age"
	"com.fha.gocan/business/core/app"
	"com.fha.gocan/business/core/commit"
	"com.fha.gocan/foundation/date"
	"com.fha.gocan/foundation/web"
	"github.com/jmoiron/sqlx"
	"net/http"
)

type handlers struct {
	Age    age.Core
	App    app.Core
	Commit commit.Core
}

func HttpMappings(connection *sqlx.DB) api.HttpMappings {
	return handlers{
		Age:    age.NewCore(connection),
		App:    app.NewCore(connection),
		Commit: commit.NewCore(connection),
	}
}

func (h handlers) GetMappings() map[string]func(w http.ResponseWriter, r *http.Request, params map[string]string) error {
	handlers := make(map[string]func(w http.ResponseWriter, r *http.Request, params map[string]string) error)
	handlers["/scenes/:sceneId/apps/:appId/code-age"] = h.queryCodeAgeHotspots
	return handlers
}

func (h handlers) queryCodeAgeHotspots(w http.ResponseWriter, r *http.Request, params map[string]string) error {
	initialDate := params["initialDate"]
	if initialDate == "" {
		initialDate = date.Today()
	}

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

	hotspots, err := h.Age.QueryCodeAgeHotspots(a.Id, initialDate, beforeTime, afterTime)
	if err != nil {
		return err
	}

	return web.Respond(w, hotspots, 200)
}
