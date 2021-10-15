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
	cr, err := c.QueryCommitRange(appId)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	before := query.Get("before")
	if before == "" {
		before = date.FormatDay(cr.MaxDate)
	}
	beforeTime, err := date.ParseDay(before)
	if err != nil {
		return time.Time{}, time.Time{}, errors.Wrap(err, "Cannot parse before parameter")
	}

	after := query.Get("after")
	if after == "" {
		after = date.FormatDay(cr.MinDate)
	}
	afterTime, err := date.ParseDay(after)
	if err != nil {
		return time.Time{}, time.Time{}, errors.Wrap(err, "Cannot parse after parameter")
	}
	return beforeTime, afterTime, nil
}

func (c Core) ExtractDateRangeFromArgs(appId string, before string, after string) (time.Time, time.Time, error) {
	cr, rangeErr := c.QueryCommitRange(appId)

	if before == "" {
		if rangeErr != nil {
			return time.Time{}, time.Time{}, errors.Wrap(rangeErr, "Commit range cannot be retrieved")
		}
		before = date.FormatDay(cr.MaxDate)
	}
	beforeTime, err := date.ParseDay(before)
	if err != nil {
		return time.Time{}, time.Time{}, errors.Wrap(err, "Invalid before date")
	}

	if after == "" {
		if rangeErr != nil {
			return time.Time{}, time.Time{}, errors.Wrap(rangeErr, "Commit range cannot be retrieved")
		}
		after = date.FormatDay(cr.MinDate)
	}
	afterTime, err := date.ParseDay(after)
	if err != nil {
		return time.Time{}, time.Time{}, errors.Wrap(err, "Invalid after date")
	}

	return beforeTime, afterTime, nil
}
