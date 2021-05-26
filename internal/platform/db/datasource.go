package db

import (
	"com.fha.gocan/internal/platform/terminal"
	"fmt"
	"github.com/jmoiron/sqlx"
	"os"
)

type SqlxDataSource struct {
	Dsn string
}

func (ds SqlxDataSource) GetConnection(ui terminal.UI) *sqlx.DB {
	ui.Say("Connecting to the database...")
	db, err := sqlx.Connect("postgres", ds.Dsn)
	if err != nil {
		ui.Failed(fmt.Sprintf("Cannot connect to the database: %v\n", err))
		os.Exit(3)
	}
	ui.Ok()
	return db
}
