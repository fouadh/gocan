package db

import (
	"com.fha.gocan/internal/terminal"
	"fmt"
	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/jmoiron/sqlx"
	"os"
)

type EmbeddedDatabase struct {
	database *embeddedpostgres.EmbeddedPostgres
}

func (ed *EmbeddedDatabase) Start(ui terminal.UI) {
	ui.Say("Starting the database...")
	ed.database = embeddedpostgres.NewDatabase()

	if err := ed.database.Start(); err != nil {
		ui.Failed(fmt.Sprintf("Cannot start the database: %+v\n", err))
		os.Exit(1)
	}
	ui.Ok()

}

func (ed *EmbeddedDatabase) Stop(ui terminal.UI) {
	ui.Say("Stopping the database...")
	if err := ed.database.Stop(); err != nil {
		ui.Failed(fmt.Sprintf("Cannot stop the database: %+v\n", err))
		os.Exit(2)
	}
	ui.Ok()
}

func (ed EmbeddedDatabase) Connect(dsn string, ui terminal.UI) *sqlx.DB {
	ui.Say("Connecting to the database...")
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		ui.Failed(fmt.Sprintf("%s Cannot connect to the database: %+v\n", err))
		os.Exit(3)
	}
	return db
}

