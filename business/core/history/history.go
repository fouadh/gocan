package history

import (
	"com.fha.gocan/business/data/store/app"
	"com.fha.gocan/business/data/store/cloc"
	"com.fha.gocan/business/data/store/commit"
	"com.fha.gocan/business/data/store/coupling"
	"com.fha.gocan/business/data/store/scene"
	"com.fha.gocan/business/data/store/stat"
	"com.fha.gocan/business/sys/git"
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

	couplings := CalculateCouplings(stats)
	c.coupling.ImportCoupling(appId, couplings)

	if err = c.stat.BulkImport(appId, stats); err != nil {
		return errors.Wrap(err, "Unable to save stats")
	}

	if err = c.cloc.ImportCloc(appId, path, commits); err != nil {
		return errors.Wrap(err, "Unable to save clocs")
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
	pairs := buildEntitiesPairs(stats)
	countCoupledEntities(stats, pairs)
	countRevisions(stats, pairs)
	return buildCouplings(pairs)
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

func countRevisions(stats []stat.Stat, pairs []*pair) {
	for _, p := range pairs {
		for _, s := range stats {
			if s.File == p.file1 {
				p.onFile1()
			}
			if s.File == p.file2 {
				p.onFile2()

			}
		}
	}
}

func countCoupledEntities(stats []stat.Stat, pairs []*pair) {
	commits := organizeEntitiesPerCommit(stats)
	for _, files := range commits {
		for _, p := range pairs {
			if _, ok := files[p.file1]; ok {
				if _, ok := files[p.file2]; ok {
					p.onCoupling()
				}
			}
		}
	}
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

func buildEntitiesPairs(stats []stat.Stat) []*pair {
	pairsMap := make(map[pair](bool))
	for _, s1 := range stats {
		for _, s2 := range stats {
			if s1.File != s2.File {
				p1 := pair{
					file1: s1.File,
					file2: s2.File,
				}
				p2 := pair{
					file1: s2.File,
					file2: s1.File,
				}
				if _, ok := pairsMap[p1]; !ok {
					if _, ok := pairsMap[p2]; !ok {
						pairsMap[p1] = true
					}
				}
			}
		}
	}

	pairs := make([]*pair, len(pairsMap))
	i := 0
	for p := range pairsMap {
		pairs[i] = &pair{file1: p.file1, file2: p.file2}
		i++
	}
	return pairs
}
