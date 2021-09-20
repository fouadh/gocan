package developer

import (
	"com.fha.gocan/business/core/app"
	"com.fha.gocan/business/core/developer"
	developer2 "com.fha.gocan/business/data/store/developer"
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

func (h *Handlers) QueryDevelopers(w http.ResponseWriter, r *http.Request, params map[string]string) error {
	appId := params["appId"]

	beforeTime, err := date.ParseDay(date.Today())
	if err != nil {
		return err
	}

	afterTime, err := date.ParseDay(date.LongTimeAgo())
	if err != nil {
		return err
	}

	devs, err := h.Developer.QueryDevelopers(appId, beforeTime, afterTime)
	if err != nil {
		return err
	}

	payload := struct {
		Developers []developer2.Developer `json:"authors"`
	}{
		Developers: devs,
	}

	return web.Respond(w, payload, 200)
}