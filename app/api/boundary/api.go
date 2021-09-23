package boundary

import (
	"com.fha.gocan/business/core/boundary"
	boundary2 "com.fha.gocan/business/data/store/boundary"
	"com.fha.gocan/foundation/web"
	"github.com/pkg/errors"
	"net/http"
)

type Handlers struct {
	Boundary boundary.Core
}

func (h *Handlers) QueryByAppId(w http.ResponseWriter, r *http.Request, params map[string]string) error  {
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