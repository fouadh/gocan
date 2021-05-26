package db

import (
	"com.fha.gocan/internal/platform/terminal"
	"fmt"
	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"os"
)

type EmbeddedDatabase struct {
	database *embeddedpostgres.EmbeddedPostgres
}

func (ed *EmbeddedDatabase) Init() {
	ed.database = embeddedpostgres.NewDatabase()
}

func (ed *EmbeddedDatabase) Start(ui terminal.UI) {
	ui.Say("Starting the embedded database...")
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