package coupling

import (
	"com.fha.gocan/business/core/app"
	"com.fha.gocan/business/core/commit"
	"com.fha.gocan/business/core/coupling"
	"com.fha.gocan/foundation/web"
	"github.com/jmoiron/sqlx"
	"net/http"
)

type Handlers struct {
	Coupling coupling.Core
	App      app.Core
	Commit   commit.Core
}

func NewHandlers(connection *sqlx.DB) Handlers {
	return Handlers{
		Coupling: coupling.NewCore(connection),
		App:      app.NewCore(connection),
		Commit:   commit.NewCore(connection),
	}
}

func (h *Handlers) BuildCouplingHierarchy(w http.ResponseWriter, r *http.Request, params map[string]string) error {
	appId := params["appId"]

	query := r.URL.Query()
	_, _, err := h.Commit.ExtractDateRangeFromQueryParams(appId, query)
	if err != nil {
		return err
	}

	a, err := h.App.QueryById(appId)
	if err != nil {
		return err
	}

	c, err := h.Coupling.BuildCouplingHierarchy(a, 0, 3)

	if err != nil {
		return err
	}

	return web.Respond(w, c, 200)
}
