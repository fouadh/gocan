package modus_operandi

import (
	modus_operandi "com.fha.gocan/business/data/store/modus-operandi"
	"github.com/jmoiron/sqlx"
	"time"
)

type Core struct {
	modusOperandi modus_operandi.Store
}

func (c Core) Query(appId string, before time.Time, after time.Time) ([]modus_operandi.WordCount, error){
	return c.modusOperandi.Query(appId, before, after)
}

func NewCore(connection *sqlx.DB) Core {
	return Core{
		modusOperandi: modus_operandi.NewStore(connection),
	}
}