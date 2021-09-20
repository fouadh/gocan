package developer

import (
	"com.fha.gocan/business/core/app"
	"com.fha.gocan/business/core/developer"
	"com.fha.gocan/foundation/date"
	"com.fha.gocan/foundation/web"
	"net/http"
)

type Handlers struct {
	App app.Core
	Developer developer.Core
}

func (h *Handlers) BuildKnowledgeMap(w http.ResponseWriter, r *http.Request, params map[string]string) error {
	appId := params["appId"]

	a, err := h.App.QueryById(appId)
	if err != nil {
		return err
	}

	beforeTime, err := date.ParseDay(date.Today())
	if err != nil {
		return err
	}

	afterTime, err := date.ParseDay(date.LongTimeAgo())
	if err != nil {
		return err
	}

	km, err := h.Developer.BuildKnowledgeMap(a, beforeTime, afterTime)

	return web.Respond(w, km, 200)
}
