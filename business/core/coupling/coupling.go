package coupling

import (
	"com.fha.gocan/business/data/store/app"
	"com.fha.gocan/business/data/store/commit"
	"com.fha.gocan/business/data/store/coupling"
	"com.fha.gocan/business/data/store/revision"
	"com.fha.gocan/business/data/store/scene"
	"com.fha.gocan/business/data/store/stat"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"sort"
	"strings"
	"time"
)

type Core struct {
	scene    scene.Store
	app      app.Store
	coupling coupling.Store
	stat     stat.Store
	commit   commit.Store
	revision revision.Store
}

func NewCore(connection *sqlx.DB) Core {
	return Core{
		scene:    scene.NewStore(connection),
		app:      app.NewStore(connection),
		coupling: coupling.NewStore(connection),
		stat:     stat.NewStore(connection),
		commit:   commit.NewStore(connection),
		revision: revision.NewStore(connection),
	}
}

func (c Core) Query(appId string, minimalCoupling float64, minimalRevisionsAverage int, beforeTime time.Time, afterTime time.Time) ([]coupling.Coupling, error) {
	stats, err := c.stat.Query(appId, beforeTime, afterTime)
	if err != nil {
		return []coupling.Coupling{}, err
	}

	couplings := CalculateCouplings(stats, minimalCoupling, float64(minimalRevisionsAverage))
	return couplings, nil
}

func (c Core) QuerySoc(appId string, before time.Time, after time.Time) ([]coupling.Soc, error) {
	return c.coupling.QuerySoc(appId, before, after)
}

func (c Core) BuildCouplingHierarchy(a app.App, minimalCoupling float64, minimalRevisionsAverage int, beforeTime time.Time, afterTime time.Time) (coupling.CouplingHierarchy, error) {
	couplings, err := c.Query(a.Id, minimalCoupling, minimalRevisionsAverage, beforeTime, afterTime)
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

func (c Core) BuildEntityCouplingHierarchy(a app.App, entity string, minimalCoupling float64, minimalRevisionsAverage int, before time.Time, after time.Time) (revision.HotspotHierarchy, error) {
	revisions, err := c.revision.QueryByAppIdAndDateRange(a.Id, before, after)
	if err != nil {
		return revision.HotspotHierarchy{}, errors.Wrap(err, "Unable to fetch revisions")
	}

	couplings, err := c.Query(a.Id, minimalCoupling, minimalRevisionsAverage, before, after)
	if err != nil {
		return revision.HotspotHierarchy{}, errors.Wrap(err, "Unable to fetch couplings")
	}

	couplingsMap := make(map[string](float64))
	maxDegree := 0.
	for _, c := range couplings {
		if c.Entity == entity {
			couplingsMap[c.Coupled] = c.Degree
		} else if c.Coupled == entity {
			couplingsMap[c.Entity] = c.Degree
		}
		if c.Degree > maxDegree {
			maxDegree = c.Degree
		}
	}

	root := revision.HotspotHierarchy{
		Name:     a.Name,
		Children: []*revision.HotspotHierarchy{},
	}

	for _, revision := range revisions {
		path := strings.Split(revision.Entity, "/")
		buildEntityCouplingNode(path, &root, revision, couplingsMap[revision.Entity] / maxDegree)
	}

	return root, nil
}

func buildEntityCouplingNode(path []string, parent *revision.HotspotHierarchy, rev revision.Revision, couplingDegree float64) *revision.HotspotHierarchy {
	existingNode := findEntityCouplingNode(parent.Children, path[0])
	if existingNode != nil {
		return buildEntityCouplingNode(path[1:], existingNode, rev, couplingDegree)
	}
	newNode := &revision.HotspotHierarchy{
		Name: path[0],
	}
	parent.Children = append(parent.Children, newNode)
	if len(path) <= 1 {
		newNode.Size = rev.Code
		newNode.Weight = couplingDegree
		return nil
	} else {
		return buildEntityCouplingNode(path[1:], parent.Children[len(parent.Children)-1], rev, couplingDegree)
	}
}

func findEntityCouplingNode(nodes []*revision.HotspotHierarchy, name string) *revision.HotspotHierarchy {
	for _, n := range nodes {
		if n.Name == name {
			return n
		}
	}
	return nil
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

type pair struct {
	file1     string
	file1Revs int
	file2     string
	file2Revs int
	count     int
}

func (p *pair) onFile1() {
	p.file1Revs++
}

func (p *pair) onFile2() {
	p.file2Revs++
}

func (p *pair) onCoupling() {
	p.count++
}

func CalculateCouplings(stats []stat.Stat, minimalCoupling float64, average float64) []coupling.Coupling {
	pairs := calculateCouplingStats(stats)
	couplings := buildCouplings(pairs, minimalCoupling, average)
	sort.Slice(couplings, func(i, j int) bool {
		return couplings[i].Degree > couplings[j].Degree
	})
	return couplings
}

func calculateCouplingStats(stats []stat.Stat) []*pair {
	commits := organizeEntitiesPerCommit(stats)
	pairsMap := make(map[pair](info))
	revisions := make(map[string](int))

	for _, files := range commits {
		for file1, _ := range files {
			revisions[file1]++
			for file2, _ := range files {
				if file1 != file2 {
					p1 := pair{
						file1: file1,
						file2: file2,
					}
					p2 := pair{
						file1: file2,
						file2: file1,
					}
					if _, ok := pairsMap[p1]; !ok {
						if _, ok := pairsMap[p2]; !ok {
							pairsMap[p1] = info{
								count: 1,
							}
						}
					} else {
						i := pairsMap[p1]
						i.count++
						pairsMap[p1] = i

					}
				}
			}
		}
	}

	pairs := make([]*pair, len(pairsMap))
	index := 0
	for p, i := range pairsMap {
		pairs[index] = &pair{file1: p.file1, file2: p.file2, count: i.count, file1Revs: revisions[p.file1], file2Revs: revisions[p.file2]}
		index++
	}
	return pairs
}

func buildCouplings(pairs []*pair, minimalCoupling float64, minimalAverage float64) []coupling.Coupling {
	couplings := []coupling.Coupling{}
	for _, p := range pairs {
		if p.count > 0 {
			average := float64(p.file1Revs+p.file2Revs) / 2
			degree := float64(p.count) / average

			if degree >= minimalCoupling && average >= minimalAverage {
				c := coupling.Coupling{
					Entity:           p.file1,
					Coupled:          p.file2,
					Degree:           degree,
					AverageRevisions: average,
				}

				couplings = append(couplings, c)
			}
		}
	}
	return couplings
}

func organizeEntitiesPerCommit(stats []stat.Stat) map[string]map[string]bool {
	commits := make(map[string](map[string](bool)))
	for _, s := range stats {
		if _, ok := commits[s.CommitId]; ok {
			commits[s.CommitId][s.File] = true
		} else {
			commits[s.CommitId] = map[string]bool{s.File: true}
		}
	}
	return commits
}

type info struct {
	count     int
	file1Revs int
	file2Revs int
}
