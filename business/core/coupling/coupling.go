package coupling

import (
	"com.fha.gocan/business/data/store/app"
	"com.fha.gocan/business/data/store/coupling"
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
	coupling coupling.Store
}

func NewCore(connection *sqlx.DB) Core {
	return Core{
		scene:    scene.NewStore(connection),
		app:      app.NewStore(connection),
		coupling: coupling.NewStore(connection),
	}
}

func (c Core) Query(appId string, minimalCoupling float64, minimalRevisionsAverage int) ([]coupling.Coupling, error) {
	return c.coupling.Query(appId, minimalCoupling, minimalRevisionsAverage)
}

func (c Core) QuerySoc(appId string, before time.Time, after time.Time) ([]coupling.Soc, error) {
	return c.coupling.QuerySoc(appId, before, after)
}

func (c Core) BuildCouplingHierarchy(a app.App, minimalCoupling float64, minimalRevisionsAverage int) (coupling.CouplingHierarchy, error) {
	couplings, err := c.Query(a.Id, minimalCoupling, minimalRevisionsAverage)
	if err != nil {
		return coupling.CouplingHierarchy{}, errors.Wrap(err, "Unable to fetch couplings")
	}

	root := coupling.CouplingHierarchy{
		Name:     "root",
		Children: []*coupling.CouplingHierarchy{},
	}

	for _, c := range couplings {
		path := strings.Split(c.Entity, "/")
		buildNode(path, &root, c)
	}

	buildCouplingNodes(&root, &root)

	return root, nil
}

func buildCouplingNodes(node *coupling.CouplingHierarchy, root *coupling.CouplingHierarchy) {
	for _, n := range node.Children {
		if n.Coupling != nil {
			for _, c := range n.Coupling {
				path := strings.Split(c, "/")
				buildCouplingNode(path[1:], root)
			}
		} else {
			buildCouplingNodes(n, root)
		}
	}
}

func buildCouplingNode(path []string, parent *coupling.CouplingHierarchy) *coupling.CouplingHierarchy {
	node := findNode(parent.Children, path[0])
	if node != nil {
		if len(path) <= 1 {
			return node
		} else {
			return buildCouplingNode(path[1:], node)
		}
	} else {
		newNode := &coupling.CouplingHierarchy{
			Name: path[0],
		}
		parent.Children = append(parent.Children, newNode)
		if len(path) <= 1 {
			return newNode
		} else {
			return buildCouplingNode(path[1:], parent.Children[len(parent.Children)-1])
		}
	}
}

func findNode(nodes []*coupling.CouplingHierarchy, name string) *coupling.CouplingHierarchy {
	for _, n := range nodes {
		if n.Name == name {
			return n
		}
	}
	return nil
}

func buildNode(path []string, parent *coupling.CouplingHierarchy, c coupling.Coupling) *coupling.CouplingHierarchy {
	entityNode := findNode(parent.Children, path[0])
	if entityNode != nil {
		if len(path) <= 1 {
			entityNode.Coupling = append(entityNode.Coupling, fmt.Sprintf("root/%s", c.Coupled))
			entityNode.Relations = append(entityNode.Relations,
				coupling.CouplingRelation{
					Coupled:          fmt.Sprintf("root/%s", c.Coupled),
					Degree:           c.Degree,
					AverageRevisions: c.AverageRevisions,
				},
			)
			return entityNode
		} else {
			return buildNode(path[1:], entityNode, c)
		}
	}
	newNode := &coupling.CouplingHierarchy{
		Name: path[0],
	}
	parent.Children = append(parent.Children, newNode)
	if len(path) <= 1 {
		newNode.Coupling = []string{fmt.Sprintf("root/%s", c.Coupled)}
		newNode.Relations = []coupling.CouplingRelation{
			{
				Coupled:          fmt.Sprintf("root/%s", c.Coupled),
				Degree:           c.Degree,
				AverageRevisions: c.AverageRevisions,
			},
		}
		return nil
	} else {
		return buildNode(path[1:], parent.Children[len(parent.Children)-1], c)
	}
}


