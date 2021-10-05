package history

import (
	"com.fha.gocan/business/data/store/app"
	"com.fha.gocan/business/data/store/cloc"
	"com.fha.gocan/business/data/store/commit"
	"com.fha.gocan/business/data/store/coupling"
	"com.fha.gocan/business/data/store/scene"
	"com.fha.gocan/business/data/store/stat"
	"com.fha.gocan/business/sys/git"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"time"
)

type Core struct {
	scene    scene.Store
	app      app.Store
	commit   commit.Store
	stat     stat.Store
	cloc     cloc.Store
	coupling coupling.Store
}

func NewCore(connection *sqlx.DB) Core {
	return Core{
		scene:    scene.NewStore(connection),
		app:      app.NewStore(connection),
		commit:   commit.NewStore(connection),
		stat:     stat.NewStore(connection),
		cloc:     cloc.NewStore(connection),
		coupling: coupling.NewStore(connection),
	}
}

func (c Core) Import(appId string, path string, before time.Time, after time.Time) error {
	commits, err := git.GetCommits(path, before, after)
	if err != nil {
		return errors.Wrap(err, "Unable to retrieve commits")
	}
	if err = c.commit.BulkImport(appId, commits); err != nil {
		return errors.Wrap(err, "Unable to save commits")
	}

	commitsMap := make(map[string]commit.Commit)
	for _, ct := range commits {
		commitsMap[ct.Id] = ct
	}
	stats, err := git.GetStats(path, before, after, commitsMap)
	if err != nil {
		return err
	}

	fmt.Println("calculate couplings...it might take some time...")
	couplings := CalculateCouplings(stats)
	fmt.Println("importing couplings...")
	c.coupling.ImportCoupling(appId, couplings)
	fmt.Println("couplings imported")

	if err = c.stat.BulkImport(appId, stats); err != nil {
		return errors.Wrap(err, "Unable to save stats")
	}

	if err = c.cloc.ImportCloc(appId, path, commits); err != nil {
		return errors.Wrap(err, "Unable to save clocs")
	}

	return nil
}

func (c Core) CheckIfCanImport(path string) error {
	ok, err := git.CheckIfAllCommited(path)
	if err != nil {
		return errors.Wrap(err, "Unable to check repo status")
	}
	if !ok {
		return errors.Errorf("The directory seems to contain files that have not been commited: please stash them or commit them before running the command.")
	}
	return nil
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

func CalculateCouplings(stats []stat.Stat) []coupling.Coupling {
	pairs := calculateCouplingStats(stats)
	return buildCouplings(pairs)
}

func calculateCouplingStats(stats []stat.Stat) []*pair {
	commits := organizeEntitiesPerCommit(stats)
	pairsMap := make(map[pair](info))
	revisions := make(map[string](int))

	fmt.Println("analyzing commits...")
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
								count:     1,
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

func buildCouplings(pairs []*pair) []coupling.Coupling {
	couplings := []coupling.Coupling{}
	for _, p := range pairs {
		if p.count > 0 {
			average := float64(p.file1Revs+p.file2Revs) / 2
			degree := float64(p.count) / average

			c := coupling.Coupling{
				Entity:           p.file1,
				Coupled:          p.file2,
				Degree:           degree,
				AverageRevisions: average,
			}

			couplings = append(couplings, c)
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
	count int
	file1Revs int
	file2Revs int
}



