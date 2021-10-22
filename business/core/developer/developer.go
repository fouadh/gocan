package developer

import (
	"com.fha.gocan/business/data/store/app"
	"com.fha.gocan/business/data/store/developer"
	"com.fha.gocan/business/data/store/revision"
	"com.fha.gocan/business/data/store/scene"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"sort"
	"strings"
	"time"
)

type Core struct {
	scene    scene.Store
	app      app.Store
	developer developer.Store
	revision revision.Store
}

func (c Core) QueryMainDevelopers(appId string, before time.Time, after time.Time) ([]developer.EntityDeveloper, error) {
	return c.developer.QueryMainDevelopers(appId, before, after)
}

func (c Core) BuildKnowledgeMap(a app.App, before time.Time, after time.Time) (developer.KnowledgeMapHierarchy, error) {
	md, err := c.developer.QueryMainDevelopers(a.Id, before, after)
	if err != nil {
		return developer.KnowledgeMapHierarchy{}, errors.Wrap(err, "Unable to fetch main developers")
	}

	revs, err := c.revision.QueryByAppIdAndDateRange(a.Id, before, after)
	if err != nil {
		return developer.KnowledgeMapHierarchy{}, errors.Wrap(err, "Unable to fetch revisions")
	}

	return buildKnowledgeMap(a.Name, revs, md), nil
}

func (c Core) QueryEntityEfforts(appId string, before time.Time, after time.Time) ([]developer.EntityEffort, error) {
	return c.developer.QueryEntityEfforts(appId, before, after)
}

func (c Core) QueryDevelopers(appId string, before time.Time, after time.Time) ([]developer.Developer, error) {
	return c.developer.QueryDevelopers(appId, before, after)
}

func (c Core) RenameDeveloper(appId string, current string, new string) error {
	return c.developer.Rename(appId, current, new)
}

func (c Core) QueryEntityEffortsForEntity(appId string, entity string, before time.Time, after time.Time) ([]developer.EntityEffort, error) {
	efforts, err := c.developer.QueryEntityEfforts(appId, before, after)
	if err != nil {
		return []developer.EntityEffort{}, nil
	}

	contributions := []developer.EntityEffort{}
	for _, e := range efforts {
		if e.Entity == entity {
			contributions = append(contributions, e)
		}
	}

	sort.Slice(contributions, func(i, j int) bool {
		return contributions[i].AuthorRevisions > contributions[j].AuthorRevisions
	})

	return contributions, nil
}

func buildKnowledgeMap(appName string, revisions []revision.Revision, developers []developer.EntityDeveloper) developer.KnowledgeMapHierarchy {
	app := developer.KnowledgeMapHierarchy{
		Name:     appName,
		Children: []*developer.KnowledgeMapHierarchy{},
	}
	devMap := make(map[string]developer.EntityDeveloper)
	for _, dev := range developers {
		devMap[dev.Entity] = dev
	}

	for _, revision := range revisions {
		dev := devMap[revision.Entity]
		path := strings.Split(revision.Entity, "/")
		buildNode(path, &app, revision, dev)
	}

	return app
}

func buildNode(path []string, parent *developer.KnowledgeMapHierarchy, revision revision.Revision, dev developer.EntityDeveloper) *developer.KnowledgeMapHierarchy {
	existingNode := findNode(parent.Children, path[0])
	if existingNode != nil {
		return buildNode(path[1:], existingNode, revision, dev)
	}
	newNode := &developer.KnowledgeMapHierarchy{
		Name:     path[0],
	}
	parent.Children = append(parent.Children, newNode)
	if len(path) <= 1 {
		newNode.Size = revision.Code
		newNode.Weight = dev.Ownership
		newNode.MainDeveloper = dev.Author

		return nil
	} else {
		return buildNode(path[1:], parent.Children[len(parent.Children)-1], revision, dev)
	}
}

func findNode(nodes []*developer.KnowledgeMapHierarchy, name string) *developer.KnowledgeMapHierarchy {
	for _, n := range nodes {
		if n.Name == name {
			return n
		}
	}
	return nil
}

func NewCore(connection *sqlx.DB) Core {
	return Core{
		scene:    scene.NewStore(connection),
		app:      app.NewStore(connection),
		developer: developer.NewStore(connection),
		revision: revision.NewStore(connection),
	}
}
