package revision

import (
	"com.fha.gocan/business/data/store/app"
	"com.fha.gocan/business/data/store/boundary"
	"com.fha.gocan/business/data/store/revision"
	"com.fha.gocan/business/data/store/scene"
	"github.com/jmoiron/sqlx"
	"github.com/pborman/uuid"
	"github.com/pkg/errors"
	"strings"
	"time"
)

type Core struct {
	revision revision.Store
	scene    scene.Store
	app      app.Store
	boundary boundary.Store
}

func NewCore(connection *sqlx.DB) Core {
	return Core{
		revision: revision.NewStore(connection),
		scene:    scene.NewStore(connection),
		app:      app.NewStore(connection),
		boundary: boundary.NewStore(connection),
	}
}

func (c Core) Query(appId string, before time.Time, after time.Time) ([]revision.Revision, error) {
	return c.revision.QueryByAppIdAndDateRange(appId, before, after)
}

func (c Core) QueryByModule(appId string, mod boundary.Module, before time.Time, after time.Time) ([]revision.Revision, error) {
	revs, err := c.Query(appId, before, after)
	if err != nil {
		return nil, err
	}
	return filterRevisionsByModule(revs, mod), nil
}

func filterRevisionsByModule(revs []revision.Revision, mod boundary.Module) []revision.Revision {
	var filteredRevs []revision.Revision

	for _, rev := range revs {
		if strings.HasPrefix(rev.Entity, mod.Path) {
			filteredRevs = append(filteredRevs, rev)
		}
	}
	return filteredRevs
}

func (c Core) RevisionTrendsByName(name string, appId string) (revision.RevisionTrends, error) {
	return c.revision.QueryTrendsByName(name, appId)
}

func (c Core) RevisionTrendsById(trendId string) ([]revision.RevisionTrend, error) {
	return c.revision.QueryTrends(trendId)
}

func (c Core) CreateRevisionTrends(name string, appId string, b boundary.Boundary, before time.Time, after time.Time) error {
	daysInRange := before.Sub(after).Hours() / 24

	trendId := uuid.New()
	entries := []revision.NewRevisionTrend{}
	for i := 0; i <= int(daysInRange); i++ {
		day := after.AddDate(0, 0, i)
		dayRevs, err := c.revision.QueryByBoundary(appId, b, day, after)
		if err != nil {
			return errors.Wrap(err, "Unable to get revisions")
		}

		entryId := uuid.New()
		trendRevs := make([]revision.TrendRevision, len(dayRevs))
		for i, rev := range dayRevs {
			trendRevs[i] = revision.TrendRevision{
				EntryId:           entryId,
				Entity:            rev.Entity,
				NumberOfRevisions: rev.NumberOfRevisions,
			}
		}

		entry := revision.NewRevisionTrend{
			Id:              entryId,
			RevisionTrendId: trendId,
			Date:            day.Format("2006-01-02"),
			Revisions:       trendRevs,
		}

		entries = append(entries, entry)
	}

	trend := revision.NewRevisionTrends{
		Id:         trendId,
		Name:       name,
		BoundaryId: b.Id,
		AppId:      appId,
		Entries:    entries,
	}

	if err := c.revision.CreateTrend(trend); err != nil {
		return errors.Wrap(err, "Unable to create trend in database")
	}

	return nil
}

func (c Core) RevisionTrendsByAppId(appId string) ([]revision.RevisionTrends, error) {
	return c.revision.QueryTrendsByAppId(appId)
}

func (c Core) QueryByBoundary(appId string, b boundary.Boundary, before time.Time, after time.Time) ([]revision.Revision, error) {
	return c.revision.QueryByBoundary(appId, b, before, after)
}
