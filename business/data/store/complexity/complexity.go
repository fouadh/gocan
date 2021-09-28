package complexity

import "github.com/jmoiron/sqlx"

type Store struct {
	connection *sqlx.DB
}

func NewStore(connection *sqlx.DB) Store {
	return Store{connection: connection}
}
