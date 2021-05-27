package db

import (
	"com.fha.gocan/internal/platform/terminal"
	"fmt"
	"github.com/jmoiron/sqlx"
	"os"
)

type DataSource interface {
	GetConnection() *sqlx.DB
}

type SqlxDataSource struct {
	Dsn string
	Ui terminal.UI
}

func (ds *SqlxDataSource) GetConnection() *sqlx.DB {
	ds.Ui.Say("Connecting to the database...")
	db, err := sqlx.Connect("postgres", ds.Dsn)
	if err != nil {
		ds.Ui.Failed(fmt.Sprintf("Cannot connect to the database: %v\n", err))
		os.Exit(3)
	}
	ds.Ui.Ok()
	return db
}
