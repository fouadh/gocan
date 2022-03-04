package boundary

import (
	"com.fha.gocan/business/api"
	"com.fha.gocan/business/core/boundary"
	boundary2 "com.fha.gocan/business/data/store/boundary"
	"com.fha.gocan/foundation/web"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"net/http"
)

type handlers struct {
	Boundary boundary.Core
}

func HttpMappings(connection *sqlx.DB) api.HttpMappings {
	return handlers{
		Boundary: boundary.NewCore(connection),
	}
}

func (h *handlers) queryByAppId(w http.ResponseWriter, r *http.Request, params map[string]string) error {
	appId := params["appId"]

	boundaries, err := h.Boundary.QueryByAppId(appId)
	if err != nil {
		return errors.Wrap(err, "Cannot retrieve boundaries")
	}

	payload := struct {
		Boundaries []boundary2.Boundary `json:"boundaries"`
	}{
		Boundaries: boundaries,
	}

	return web.Respond(w, payload, 200)
}

func (h handlers) GetMappings() map[string]func(w http.ResponseWriter, r *http.Request, params map[string]string) error {
	handlers := make(map[string]func(w http.ResponseWriter, r *http.Request, params map[string]string) error)
	handlers["/scenes/:sceneId/apps/:appId/boundaries"] = h.queryByAppId
	return handlers
}
