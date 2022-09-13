package age

import (
	"com.fha.gocan/business/data/store/age"
	"com.fha.gocan/business/data/store/revision"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"strings"
	"time"
)

type Core struct {
	age      age.Store
	revision revision.Store
}

func NewCore(connection *sqlx.DB) Core {
	return Core{
		age:      age.NewStore(connection),
		revision: revision.NewStore(connection),
	}
}

func (c Core) GetCodeAge(appId string, initialDate string, before time.Time, after time.Time) ([]age.EntityAge, error) {
	return c.age.QueryEntityAge(appId, initialDate, before, after)
}

func (c Core) QueryCodeAgeHotspots(appId string, initialDate string, before time.Time, after time.Time) (age.EntityAgeHierarchy, error) {
	data, err := c.GetCodeAge(appId, initialDate, before, after)
	if err != nil {
		return age.EntityAgeHierarchy{}, errors.Wrap(err, "Unable to fetch code age")
	}

	revs, err := c.revision.QueryByAppIdAndDateRange(appId, before, after)
	if err != nil {
		return age.EntityAgeHierarchy{}, errors.Wrap(err, "Unable to fetch revisions")
	}

	return buildHotspots(data, revs), nil
}

func buildHotspots(data []age.EntityAge, revs []revision.Revision) age.EntityAgeHierarchy {
	root := age.EntityAgeHierarchy{
		Name:     "root",
		Children: []*age.EntityAgeHierarchy{},
	}

	revMap := make(map[string]revision.Revision)
	for _, rev := range revs {
		revMap[rev.Entity] = rev
	}

	var maxAge = 0
	for _, ea := range data {
		if maxAge < ea.Age {
			maxAge = ea.Age
		}
	}

	for _, ea := range data {
		path := strings.Split(ea.Name, "/")
		buildNode(path, &root, ea, revMap, maxAge)
	}

	return root
}

func buildNode(path []string, parent *age.EntityAgeHierarchy, ea age.EntityAge, revMap map[string]revision.Revision, maxAge int) *age.EntityAgeHierarchy {
	existingNode := findNode(parent.Children, path[0])
	if existingNode != nil {
		return buildNode(path[1:], existingNode, ea, revMap, maxAge)
	}
	newNode := &age.EntityAgeHierarchy{
		Name: path[0],
	}
	parent.Children = append(parent.Children, newNode)
	if len(path) <= 1 {
		newNode.Size = revMap[ea.Name].Code
		newNode.Weight = float64(ea.Age) / float64(maxAge)
		return nil
	} else {
		return buildNode(path[1:], parent.Children[len(parent.Children)-1], ea, revMap, maxAge)
	}
}

func findNode(nodes []*age.EntityAgeHierarchy, name string) *age.EntityAgeHierarchy {
	for _, n := range nodes {
		if n.Name == name {
			return n
		}
	}
	return nil
}
