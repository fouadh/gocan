package coupling

import (
	"com.fha.gocan/business/core/app"
	"com.fha.gocan/business/core/coupling"
	"com.fha.gocan/foundation/date"
	"com.fha.gocan/foundation/web"
	"net/http"
)

type Handlers struct {
	Coupling coupling.Core
	App app.Core
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

	a, err := h.App.QueryById(appId)
	if err != nil {
		return err
	}

	c, err := h.Coupling.BuildCouplingHierarchy(a, beforeTime, afterTime, 0, 3)

	if err != nil {
		return err
	}

	return web.Respond(w, c, 200)
}
