package modus_operandi

import (
	modus_operandi "com.fha.gocan/business/core/modus-operandi"
	modus_operandi2 "com.fha.gocan/business/data/store/modus-operandi"
	"com.fha.gocan/foundation/date"
	"com.fha.gocan/foundation/web"
	"github.com/jmoiron/sqlx"
	"net/http"
)

type Handlers struct {
	ModusOperandi modus_operandi.Core
}

func NewHandlers(connection *sqlx.DB) Handlers {
	return Handlers{
		ModusOperandi: modus_operandi.NewCore(connection),
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

	mo, err := h.ModusOperandi.Query(appId, beforeTime, afterTime)
	if err != nil {
		return err
	}

	payload := struct {
		ModusOperandi []modus_operandi2.WordCount `json:"modusOperandi"`
	}{
		ModusOperandi: mo,
	}

	return web.Respond(w, payload, 200)
}
