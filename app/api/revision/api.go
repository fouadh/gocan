package revision

import (
	"com.fha.gocan/business/core/revision"
	revision2 "com.fha.gocan/business/data/store/revision"
	"com.fha.gocan/foundation/date"
	"com.fha.gocan/foundation/web"
	"net/http"
)

type Handlers struct {
	Revision   revision.Core
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

	revs, err := h.Revision.GetRevisions(appId, beforeTime, afterTime)
	if err != nil {
		return err
	}

	result := struct {
		Revisions []revision2.Revision `json:"revisions"`
	}{
		Revisions: revs,
	}

	return web.Respond(w, result, 200)
}
