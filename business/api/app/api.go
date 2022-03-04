package app

import (
	"com.fha.gocan/business/api"
	"com.fha.gocan/business/core/app"
	"com.fha.gocan/business/core/commit"
	app2 "com.fha.gocan/business/data/store/app"
	"com.fha.gocan/foundation/date"
	"com.fha.gocan/foundation/web"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"net/http"
)

type handlers struct {
	App    app.Core
	Commit commit.Core
}

func HttpMappings(connection *sqlx.DB) api.HttpMappings {
	return handlers{
		App:    app.NewCore(connection),
		Commit: commit.NewCore(connection),
	}
}

func (h *handlers) queryAll(w http.ResponseWriter, r *http.Request, params map[string]string) error {
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
			MinDate: date.FormatDay(cr.MinDay()),
			MaxDate: date.FormatDay(cr.MaxDay()),
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

func (h *handlers) queryById(w http.ResponseWriter, r *http.Request, params map[string]string) error {
	appId := params["appId"]
	a, err := h.App.QueryById(appId)
	if err != nil {
		return err
	}

	return web.Respond(w, a, 200)
}

func (h *handlers) queryEntities(w http.ResponseWriter, r *http.Request, params map[string]string) error {
	appId := params["appId"]

	entities, err := h.App.QueryEntities(appId)
	if err != nil {
		return errors.Wrap(err, "Unable to get entities")
	}

	files := make([]string, len(entities))
	for i, e := range entities {
		files[i] = e.Name
	}

	payload := struct {
		Entities []string `json:"entities"`
	}{
		Entities: files,
	}

	return web.Respond(w, payload, 200)
}

func (h handlers) GetMappings() map[string]func(w http.ResponseWriter, r *http.Request, params map[string]string) error {
	handlers := make(map[string]func(w http.ResponseWriter, r *http.Request, params map[string]string) error)
	handlers["/scenes/:sceneId/apps"] = h.queryAll
	handlers["/scenes/:sceneId/apps/:appId"] = h.queryById
	handlers["/scenes/:sceneId/apps/:appId/entities"] = h.queryEntities
	return handlers
}
