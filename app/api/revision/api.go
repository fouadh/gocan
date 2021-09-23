package revision

import (
	"com.fha.gocan/business/core/app"
	"com.fha.gocan/business/core/boundary"
	"com.fha.gocan/business/core/commit"
	"com.fha.gocan/business/core/revision"
	revision2 "com.fha.gocan/business/data/store/revision"
	"com.fha.gocan/foundation/date"
	"com.fha.gocan/foundation/web"
	"github.com/pkg/errors"
	"net/http"
)

type Handlers struct {
	Revision revision.Core
	App      app.Core
	Boundary boundary.Core
	Commit   commit.Core
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

	revs, err := h.Revision.Query(appId, beforeTime, afterTime)
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

func (h *Handlers) QueryHotspots(w http.ResponseWriter, r *http.Request, params map[string]string) error {
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

	hotspots, err := h.Revision.QueryHotspots(a, beforeTime, afterTime)

	payload := struct {
		Name     string                       `json:"name"`
		Children []revision2.HotspotHierarchy `json:"children"`
	}{
		Name:     "root",
		Children: []revision2.HotspotHierarchy{hotspots},
	}

	return web.Respond(w, payload, 200)
}

func (h *Handlers) QueryRevisionsTrends(w http.ResponseWriter, r *http.Request, params map[string]string) error {
	appId := params["appId"]
	a, err := h.App.QueryById(appId)
	if err != nil {
		return err
	}

	query := r.URL.Query()
	boundaryId := query.Get("boundaryId")
	b, err := h.Boundary.QueryByBoundaryId(boundaryId)
	if err != nil {
		return errors.Wrap(err, "Boundary not found")
	}

	cr, err := h.Commit.QueryCommitRange(appId)
	if err != nil {
		return err
	}

	before := query.Get("before")
	if before == "" {
		before = date.FormatDay(cr.MaxDate)
	}
	beforeTime, err := date.ParseDay(before)
	if err != nil {
		return errors.Wrap(err, "Cannot parse before parameter")
	}

	after := query.Get("after")
	if after == "" {
		after = date.FormatDay(cr.MinDate)
	}
	afterTime, err := date.ParseDay(after)
	if err != nil {
		return errors.Wrap(err, "Cannot parse after parameter")
	}

	trends, err := h.Revision.RevisionTrends(a.Id, b, beforeTime, afterTime)
	if err != nil {
		return err
	}

	payload := struct {
		Trends []revision2.RevisionTrend `json:"trends"`
	}{
		Trends: trends,
	}

	return web.Respond(w, payload, 200)
}
