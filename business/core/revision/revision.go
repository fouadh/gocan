package revision

import (
	"com.fha.gocan/business/data/store/app"
	"com.fha.gocan/business/data/store/revision"
	"com.fha.gocan/business/data/store/scene"
	context "com.fha.gocan/foundation"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"strings"
	"time"
)

type Core struct {
	revision revision.Store
	scene scene.Store
	app   app.Store
}

func NewCore(connection *sqlx.DB) Core {
	return Core{
		revision: revision.NewStore(connection),
		scene: scene.NewStore(connection),
		app: app.NewStore(connection),
	}
}

func (c Core) GetRevisions(ctx context.Context, appName string, sceneName string, before time.Time, after time.Time) ([]revision.Revision, error) {
	s, err := c.scene.QueryByName(sceneName)
	if err != nil {
		return []revision.Revision{}, errors.Wrap(err, "Scene not found")
	}

	a, err := c.app.QueryBySceneIdAndName(s.Id, appName)
	if err != nil {
		return []revision.Revision{}, errors.Wrap(err, "App not found")
	}

	return c.revision.QueryByAppIdAndDateRange(a.Id, before, after)
}

func (c Core) GetHotspots(ctx context.Context, appName string, sceneName string, before time.Time, after time.Time) (revision.HotspotHierarchy, error) {
	s, err := c.scene.QueryByName(sceneName)
	if err != nil {
		return revision.HotspotHierarchy{}, errors.Wrap(err, "Scene not found")
	}

	a, err := c.app.QueryBySceneIdAndName(s.Id, appName)
	if err != nil {
		return revision.HotspotHierarchy{}, errors.Wrap(err, "App not found")
	}

	revs, err := c.revision.QueryByAppIdAndDateRange(a.Id, before, after)
	if err != nil {
		return revision.HotspotHierarchy{}, errors.Wrap(err, "Unable to fetch revisions")
	}

	return buildHotspots(appName, revs), nil
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
		Name:     path[0],
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


