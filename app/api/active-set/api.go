package active_set

import (
	active_set "com.fha.gocan/business/core/active-set"
	active_set2 "com.fha.gocan/business/data/store/active-set"
	"com.fha.gocan/foundation/date"
	"com.fha.gocan/foundation/web"
	"net/http"
)

type Handlers struct {
	ActiveSet active_set.Core
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