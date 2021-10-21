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

func (c Core) QueryHotspots(a app.App, before time.Time, after time.Time) (revision.HotspotHierarchy, error) {
	revs, err := c.revision.QueryByAppIdAndDateRange(a.Id, before, after)
	if err != nil {
		return revision.HotspotHierarchy{}, errors.Wrap(err, "Unable to fetch revisions")
	}

	return buildHotspots(a.Name, revs), nil
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
		AppId: appId,
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

func buildHotspots(appName string, revisions []revision.Revision) revision.HotspotHierarchy {
	root := revision.HotspotHierarchy{
		Name:     appName,
		Children: []*revision.HotspotHierarchy{},
	}

	for _, revision := range revisions {
		path := strings.Split(revision.Entity, "/")
		buildNode(path, &root, revision)
	}

	return root
}

func buildNode(path []string, parent *revision.HotspotHierarchy, rev revision.Revision) *revision.HotspotHierarchy {
	existingNode := findNode(parent.Children, path[0])
	if existingNode != nil {
		return buildNode(path[1:], existingNode, rev)
	}
	newNode := &revision.HotspotHierarchy{
		Name: path[0],
	}
	parent.Children = append(parent.Children, newNode)
	if len(path) <= 1 {
		newNode.Size = rev.Code
		newNode.Weight = rev.NormalizedNumberOfRevisions
		return nil
	} else {
		return buildNode(path[1:], parent.Children[len(parent.Children)-1], rev)
	}
}

func findNode(nodes []*revision.HotspotHierarchy, name string) *revision.HotspotHierarchy {
	for _, n := range nodes {
		if n.Name == name {
			return n
		}
	}
	return nil
}
