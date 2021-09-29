package complexity

import (
	"com.fha.gocan/business/core/complexity"
	complexity2 "com.fha.gocan/business/data/store/complexity"
	"com.fha.gocan/foundation/web"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"net/http"
)

type Handlers struct {
	Complexity complexity.Core
}

func NewHandlers(connection *sqlx.DB) Handlers {
	return Handlers{
		Complexity: complexity.NewCore(connection),
	}
}

func (h *Handlers) QueryAnalyses(w http.ResponseWriter, r *http.Request, params map[string]string) error {
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