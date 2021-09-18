package active_set

import (
	active_set "com.fha.gocan/business/data/store/active-set"
	"github.com/jmoiron/sqlx"
	"sort"
	"time"
)

type Core struct {
	activeSet active_set.Store
}

func (c Core) Query(appId string, before time.Time, after time.Time) ([]active_set.ActiveSet, error) {
	opened, err := c.activeSet.QueryOpenedEntities(appId, before, after)
	if err != nil {
		return []active_set.ActiveSet{}, err
	}

	closed, err := c.activeSet.QueryClosedEntities(appId, before, after)
	if err != nil {
		return []active_set.ActiveSet{}, err
	}

	mergedResults := make(map[time.Time]active_set.ActiveSet)

	for _, each := range opened {
		item := active_set.ActiveSet{
			Date:   each.Date,
			Opened: each.Count,
		}
		mergedResults[item.Date] = item
	}

	for _, each := range closed {
		item := mergedResults[each.Date]
		item.Closed = each.Count
		mergedResults[each.Date] = item
	}

	results := make([]active_set.ActiveSet, 0, len(mergedResults))
	for _, item := range mergedResults {
		results = append(results, item)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Date.Before(results[j].Date)
	})

	return results, nil
}

func NewCore(connection *sqlx.DB) Core {
	return Core{
		activeSet: active_set.NewStore(connection),
	}
}
