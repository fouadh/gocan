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
