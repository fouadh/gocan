package app

import (
	"com.fha.gocan/business/core/app"
	app2 "com.fha.gocan/business/data/store/app"
	"com.fha.gocan/foundation/date"
	"com.fha.gocan/foundation/web"
	"net/http"
)

type Handlers struct {
	App   app.Core
}

func (h *Handlers) QueryAll(w http.ResponseWriter, r *http.Request, params map[string]string) error {
	sceneId := params["sceneId"]
	apps, err := h.App.QueryBySceneId(sceneId)

	if err != nil {
		return err
	}

	summaries := []app2.Summary{}

	beforeTime, err := date.ParseDay(date.Today())
	if err != nil {
		return err
	}

	afterTime, err := date.ParseDay(date.LongTimeAgo())
	if err != nil {
		return err
	}

	for _, a := range apps {
		summary, err := h.App.QuerySummary(a.Id, beforeTime, afterTime)
		if err != nil {
			return err
		}
		summaries = append(summaries, summary)
	}

	result := struct {
		Applications []app2.Summary `json:"apps"`
	}{
		Applications: summaries,
	}

	return web.Respond(w, result, 200)
}