package db

import (
	"com.fha.gocan/internal/platform/terminal"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type DataSource interface {
	GetConnection() (*sqlx.DB, error)
}

type SqlxDataSource struct {
	Dsn string
	Ui terminal.UI
}

func (ds *SqlxDataSource) GetConnection() (*sqlx.DB, error) {
	ds.Ui.Say("Connecting to the database...")
	db, err := sqlx.Connect("postgres", ds.Dsn)
	if err != nil {
		return nil, errors.Wrap(err, "Cannot connect to the database")
	}
	ds.Ui.Ok()
	return db, nil
}
