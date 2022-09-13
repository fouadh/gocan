package age

import (
	"com.fha.gocan/business/data/store/age"
	"github.com/jmoiron/sqlx"
	"time"
)

type Core struct {
	age age.Store
}

func (c Core) GetCodeAge(appId string, initialDate string, before time.Time, after time.Time) ([]age.EntityAge, error) {
	return c.age.QueryEntityAge(appId, initialDate, before, after)
}

func NewCore(connection *sqlx.DB) Core {
	return Core{
		age: age.NewStore(connection),
	}
}
