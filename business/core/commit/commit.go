package commit

import (
	"com.fha.gocan/business/data/store/commit"
	"github.com/jmoiron/sqlx"
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
