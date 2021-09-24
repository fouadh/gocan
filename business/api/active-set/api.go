package active_set

import (
	active_set "com.fha.gocan/business/core/active-set"
	"com.fha.gocan/business/core/commit"
	active_set2 "com.fha.gocan/business/data/store/active-set"
	"com.fha.gocan/foundation/web"
	"github.com/jmoiron/sqlx"
	"net/http"
)

type Handlers struct {
	ActiveSet active_set.Core
	Commit    commit.Core
}

func NewHandlers(connection *sqlx.DB) Handlers {
	return Handlers{
		ActiveSet: active_set.NewCore(connection),
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

	as, err := h.ActiveSet.Query(appId, beforeTime, afterTime)
	if err != nil {
		return err
	}

	payload := struct {
		ActiveSet []active_set2.ActiveSet `json:"activeSet"`
	}{
		ActiveSet: as,
	}

	return web.Respond(w, payload, 200)
}
