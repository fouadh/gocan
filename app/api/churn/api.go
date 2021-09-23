package churn

import (
	"com.fha.gocan/business/core/churn"
	"com.fha.gocan/business/core/commit"
	churn2 "com.fha.gocan/business/data/store/churn"
	"com.fha.gocan/foundation/web"
	"github.com/jmoiron/sqlx"
	"net/http"
)

type Handlers struct {
	Churn  churn.Core
	Commit commit.Core
}

func NewHandlers(connection *sqlx.DB) Handlers {
	return Handlers{
		Churn: churn.NewCore(connection),
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

	codeChurn, err := h.Churn.QueryCodeChurn(appId, beforeTime, afterTime)
	if err != nil {
		return err
	}

	result := struct {
		CodeChurn []churn2.CodeChurn `json:"codeChurn"`
	}{
		CodeChurn: codeChurn,
	}

	return web.Respond(w, result, 200)
}
