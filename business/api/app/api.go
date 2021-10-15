package app

import (
	"com.fha.gocan/business/core/app"
	"com.fha.gocan/business/core/commit"
	app2 "com.fha.gocan/business/data/store/app"
	"com.fha.gocan/foundation/date"
	"com.fha.gocan/foundation/web"
	"github.com/jmoiron/sqlx"
	"net/http"
)

type Handlers struct {
	App    app.Core
	Commit commit.Core
}

func NewHandlers(connection *sqlx.DB) Handlers {
	return Handlers{
		App: app.NewCore(connection),
		Commit: commit.NewCore(connection),
	}
}

func (h *Handlers) QueryAll(w http.ResponseWriter, r *http.Request, params map[string]string) error {
	sceneId := params["sceneId"]
	apps, err := h.App.QueryBySceneId(sceneId)

	if err != nil {
		return err
	}

	summaries := []app2.Summary{}
	for _, a := range apps {
		cr, err := h.Commit.QueryCommitRange(a.Id)
		if err != nil {
			return err
		}

		summary, err := h.App.QuerySummary(a.Id, cr.MaxDate, cr.MinDate)
		summary.DateRange = app2.DateRange{
			MinDate: date.FormatDay(cr.MinDate),
			MaxDate: date.FormatDay(cr.MaxDate),
		}
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

func (h *Handlers) QueryById(w http.ResponseWriter, r *http.Request, params map[string]string) error {
	appId := params["appId"]
	a, err := h.App.QueryById(appId)
	if err != nil {
		return err
	}

	return web.Respond(w, a, 200)
}
