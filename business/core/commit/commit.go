package commit

import (
	"com.fha.gocan/business/data/store/commit"
	"com.fha.gocan/foundation/date"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"net/url"
	"time"
)

type Core struct {
	commit commit.Store
}

func (c Core) QueryCommitRange(appId string) (commit.CommitRange, error) {
	return c.commit.QueryCommitRange(appId)
}

func NewCore(connection *sqlx.DB) Core {
	return Core{
		commit: commit.NewStore(connection),
	}
}

func (c Core) ExtractDateRangeFromQueryParams(appId string, query url.Values) (time.Time, time.Time, error) {
	before := query.Get("before")
	after := query.Get("after")

	return c.ExtractDateRangeFromArgs(appId, before, after)
}

func (c Core) ExtractDateRangeFromArgs(appId string, before string, after string) (time.Time, time.Time, error) {
	cr, rangeErr := c.QueryCommitRange(appId)
	minDay := cr.MinDay()
	maxDay := cr.MaxDay()

	if before == "" {
		if rangeErr != nil {
			return time.Time{}, time.Time{}, errors.Wrap(rangeErr, "Commit range cannot be retrieved")
		}
		before = date.FormatDay(maxDay)
	}
	beforeTime, err := date.ParseDay(before)
	if err != nil {
		return time.Time{}, time.Time{}, errors.Wrap(err, "Invalid before date")
	}

	if after == "" {
		if rangeErr != nil {
			return time.Time{}, time.Time{}, errors.Wrap(rangeErr, "Commit range cannot be retrieved")
		}
		after = date.FormatDay(minDay)
	}
	afterTime, err := date.ParseDay(after)
	if err != nil {
		return time.Time{}, time.Time{}, errors.Wrap(err, "Invalid after date")
	}

	if (afterTime.Before(minDay) && !afterTime.Equal(minDay)) || afterTime.After(maxDay)  {
		return time.Time{}, time.Time{}, errors.Errorf("Invalid after date: it should be between " + date.FormatDay(minDay) + " and " + date.FormatDay(maxDay) + " instead of " + date.FormatDay(afterTime))
	}

	if beforeTime.Before(minDay) || (beforeTime.After(maxDay) && !beforeTime.Equal(maxDay)) {
		return time.Time{}, time.Time{}, errors.Errorf("Invalid before date: it should be between " + date.FormatDay(minDay) + " and " + date.FormatDay(maxDay) + " instead of " + date.FormatDay(beforeTime))
	}

	if beforeTime.Before(afterTime) {
		return time.Time{}, time.Time{}, errors.Errorf("After date must be less that Before date: shall the dates be inverted ?")
	}

	return beforeTime, afterTime, nil
}