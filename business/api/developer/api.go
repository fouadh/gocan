package developer

import (
	"com.fha.gocan/business/core/app"
	"com.fha.gocan/business/core/commit"
	"com.fha.gocan/business/core/developer"
	developer2 "com.fha.gocan/business/data/store/developer"
	"com.fha.gocan/foundation/web"
	"github.com/jmoiron/sqlx"
	"net/http"
)

type Handlers struct {
	App       app.Core
	Developer developer.Core
	Commit    commit.Core
}

func NewHandlers(connection *sqlx.DB) Handlers {
	return Handlers{
		App:       app.NewCore(connection),
		Developer: developer.NewCore(connection),
		Commit:    commit.NewCore(connection),
	}
}

func (h *Handlers) BuildKnowledgeMap(w http.ResponseWriter, r *http.Request, params map[string]string) error {
	appId := params["appId"]

	a, err := h.App.QueryById(appId)
	if err != nil {
		return err
	}

	query := r.URL.Query()
	beforeTime, afterTime, err := h.Commit.ExtractDateRangeFromQueryParams(appId, query)
	if err != nil {
		return err
	}

	km, err := h.Developer.BuildKnowledgeMap(a, beforeTime, afterTime)

	payload := struct {
		Name     string                             `json:"name"`
		Children []developer2.KnowledgeMapHierarchy `json:"children"`
	}{
		Name:     "root",
		Children: []developer2.KnowledgeMapHierarchy{km},
	}

	return web.Respond(w, payload, 200)
}

func (h *Handlers) QueryDevelopers(w http.ResponseWriter, r *http.Request, params map[string]string) error {
	appId := params["appId"]

	query := r.URL.Query()
	beforeTime, afterTime, err := h.Commit.ExtractDateRangeFromQueryParams(appId, query)
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
