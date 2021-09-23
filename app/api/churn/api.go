package churn

import (
	"com.fha.gocan/business/core/churn"
	churn2 "com.fha.gocan/business/data/store/churn"
	"com.fha.gocan/foundation/date"
	"com.fha.gocan/foundation/web"
	"github.com/jmoiron/sqlx"
	"net/http"
)

type Handlers struct {
	Churn churn.Core
}

func NewHandlers(connection *sqlx.DB) Handlers {
	return Handlers{
		Churn: churn.NewCore(connection),
	}
}

func (h *Handlers) Query(w http.ResponseWriter, r *http.Request, params map[string]string) error {
	appId := params["appId"]

	beforeTime, err := date.ParseDay(date.Today())
	if err != nil {
		return err
	}

	afterTime, err := date.ParseDay(date.LongTimeAgo())
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

