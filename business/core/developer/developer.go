package developer

import (
	"com.fha.gocan/business/data/store/app"
	"com.fha.gocan/business/data/store/developer"
	"com.fha.gocan/business/data/store/revision"
	"com.fha.gocan/business/data/store/scene"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"strings"
	"time"
)

type Core struct {
	scene    scene.Store
	app      app.Store
	developer developer.Store
	revision revision.Store
}

func (c Core) QueryMainDevelopers(sceneName string, appName string, before time.Time, after time.Time) ([]developer.EntityDeveloper, error) {
	s, err := c.scene.QueryByName(sceneName)
	if err != nil {
		return []developer.EntityDeveloper{}, fmt.Errorf("unable to retrieve scene %s", sceneName)
	}

	a, err := c.app.QueryBySceneIdAndName(s.Id, appName)
	if err != nil {
		return []developer.EntityDeveloper{}, fmt.Errorf("unable to retrieve app %s linked to the scene %s", appName, sceneName)
	}

	return c.developer.QueryMainDevelopers(a.Id, before, after)
}

func (c Core) BuildKnowledgeMap(sceneName string, appName string, before time.Time, after time.Time) (developer.KnowledgeMapHierarchy, error) {
	s, err := c.scene.QueryByName(sceneName)
	if err != nil {
		return developer.KnowledgeMapHierarchy{}, fmt.Errorf("unable to retrieve scene %s", sceneName)
	}

	a, err := c.app.QueryBySceneIdAndName(s.Id, appName)
	if err != nil {
		return developer.KnowledgeMapHierarchy{}, fmt.Errorf("unable to retrieve app %s linked to the scene %s", appName, sceneName)
	}

	md, err := c.developer.QueryMainDevelopers(a.Id, before, after)
	if err != nil {
		return developer.KnowledgeMapHierarchy{}, errors.Wrap(err, "Unable to fetch main developers")
	}

	revs, err := c.revision.QueryByAppIdAndDateRange(a.Id, before, after)
	if err != nil {
		return developer.KnowledgeMapHierarchy{}, errors.Wrap(err, "Unable to fetch revisions")
	}

	return buildKnowledgeMap(appName, revs, md), nil
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
