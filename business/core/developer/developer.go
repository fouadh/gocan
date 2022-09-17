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
	scene     scene.Store
	app       app.Store
	developer developer.Store
	revision  revision.Store
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

	efforts, err := c.QueryEntityEfforts(a.Id, before, after)
	if err != nil {
		return developer.KnowledgeMapHierarchy{}, errors.Wrap(err, "Unable to fetch entity efforts")
	}

	return buildKnowledgeMap(a.Name, revs, md, efforts), nil
}

func (c Core) QueryEntityEffortsPerAuthor(appId string, before time.Time, after time.Time) ([]developer.EntityEffortPerAuthor, error) {
	return c.developer.QueryEntityEffortsPerAuthor(appId, before, after)
}

type Network struct {
	Nodes []Node `json:"nodes"`
	Links []Link `json:"links"`
}

type Node struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Link struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

func (c Core) QueryEntityEffortsGraph(appId string, before time.Time, after time.Time) (Network, error) {
	efforts, err := c.QueryEntityEffortsPerAuthor(appId, before, after)
	if err != nil {
		return Network{}, errors.Wrap(err, "Unable to query entity efforts")
	}

	devsByName := make(map[string]bool)
	entities := make(map[string][]string)

	var nodes []Node
	for _, effort := range efforts {
		if !devsByName[effort.Author] {
			devsByName[effort.Author] = true
			nodes = append(nodes, Node{
				Id:   effort.Author,
				Name: effort.Author,
			})
		}

		if entities[effort.Entity] == nil {
			entities[effort.Entity] = []string{}
		}

		entities[effort.Entity] = append(entities[effort.Entity], effort.Author)
	}

	return Network{
		Nodes: nodes,
		Links: buildLinks(entities),
	}, nil
}

// ouch ugly algorithm
func buildLinks(entities map[string][]string) []Link {
	var links []Link
	var linksByDev = make(map[string]map[string]bool)

	for _, devs := range entities {
		for _, dev1 := range devs {
			for _, dev2 := range devs {
				if dev1 != dev2 {
					if linksByDev[dev1] == nil {
						linksByDev[dev1] = make(map[string]bool)
					}
					if linksByDev[dev2] == nil {
						linksByDev[dev2] = make(map[string]bool)
					}

					if !linksByDev[dev1][dev2] && !linksByDev[dev2][dev1] {
						link := Link{
							Source: dev1,
							Target: dev2,
						}
						links = append(links, link)
						linksByDev[dev1][dev2] = true
					}
				}
			}
		}
	}
	return links
}

func (c Core) QueryEntityEfforts(appId string, before time.Time, after time.Time) ([]developer.EntityEffort, error) {
	return c.developer.QueryDevelopmentEffort(appId, before, after)
}

func (c Core) QueryDevelopers(appId string, before time.Time, after time.Time) ([]developer.Developer, error) {
	return c.developer.QueryDevelopers(appId, before, after)
}

func (c Core) RenameDeveloper(appId string, current string, new string) error {
	return c.developer.Rename(appId, current, new)
}

func (c Core) QueryEntityEffortsForEntity(appId string, entity string, before time.Time, after time.Time) ([]developer.EntityEffortPerAuthor, error) {
	efforts, err := c.developer.QueryEntityEffortsPerAuthor(appId, before, after)
	if err != nil {
		return []developer.EntityEffortPerAuthor{}, nil
	}

	var contributions []developer.EntityEffortPerAuthor
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

func (c Core) CreateTeam(appId string, teamName string, members []string) (developer.Team, error) {
	return c.developer.CreateTeam(developer.NewTeam{
		Name:    teamName,
		AppId:   appId,
		Members: members,
	})
}

func (c Core) DeleteTeam(appId string, teamName string) error {
	return c.developer.DeleteTeam(appId, teamName)
}

func buildKnowledgeMap(appName string,
	revisions []revision.Revision,
	developers []developer.EntityDeveloper,
	efforts []developer.EntityEffort) developer.KnowledgeMapHierarchy {
	app := developer.KnowledgeMapHierarchy{
		Name:     appName,
		Children: []*developer.KnowledgeMapHierarchy{},
	}
	devMap := make(map[string]developer.EntityDeveloper)
	for _, dev := range developers {
		devMap[dev.Entity] = dev
	}

	effortMap := make(map[string]developer.EntityEffort)
	for _, effort := range efforts {
		effortMap[effort.Entity] = effort
	}

	for _, revision := range revisions {
		dev := devMap[revision.Entity]
		effort := effortMap[revision.Entity]
		path := strings.Split(revision.Entity, "/")
		buildNode(path, &app, revision, dev, effort)
	}

	return app
}

func buildNode(path []string, parent *developer.KnowledgeMapHierarchy, revision revision.Revision, dev developer.EntityDeveloper, effort developer.EntityEffort) *developer.KnowledgeMapHierarchy {
	existingNode := findNode(parent.Children, path[0])
	if existingNode != nil {
		return buildNode(path[1:], existingNode, revision, dev, effort)
	}
	newNode := &developer.KnowledgeMapHierarchy{
		Name: path[0],
	}
	parent.Children = append(parent.Children, newNode)
	if len(path) <= 1 {
		newNode.Size = revision.Code
		newNode.Weight = dev.Ownership
		newNode.MainDeveloper = dev.Author
		newNode.Effort = effort.Effort
		if effort.Effort <= 0.25 {
			newNode.DevDiffusion = 0.25
		} else if effort.Effort <= 0.5 {
			newNode.DevDiffusion = 0.5
		} else if effort.Effort <= 0.75 {
			newNode.DevDiffusion = 0.75
		} else {
			newNode.DevDiffusion = 1.0
		}
		return nil
	} else {
		return buildNode(path[1:], parent.Children[len(parent.Children)-1], revision, dev, effort)
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
		scene:     scene.NewStore(connection),
		app:       app.NewStore(connection),
		developer: developer.NewStore(connection),
		revision:  revision.NewStore(connection),
	}
}
