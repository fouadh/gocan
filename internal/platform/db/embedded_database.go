package db

import (
	"com.fha.gocan/internal/platform/config"
	"com.fha.gocan/internal/platform/terminal"
	"fmt"
	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
)

type EmbeddedDatabase struct {
	database *embeddedpostgres.EmbeddedPostgres
	Config *config.Config
}

func (ed *EmbeddedDatabase) Start(ui terminal.UI) {
	ed.database = embeddedpostgres.NewDatabase(embeddedpostgres.DefaultConfig().
		Username(ed.Config.User).
		Password(ed.Config.Password).
		Database(ed.Config.Database).
		Port(uint32(ed.Config.Port)))

	ui.Say("Starting the embedded database...")
	if err := ed.database.Start(); err != nil {
		ui.Failed(fmt.Sprintf("Cannot start the database: %+v\n", err))
		os.Exit(1)
	}
	ui.Ok()

}

func (ed *EmbeddedDatabase) Stop(ui terminal.UI) {
	ui.Say("Stopping the database...")
	usr, _ := user.Current()
	dir := usr.HomeDir
	// todo the path to the db should come from the configuration
	stopPostgres(filepath.Join(dir, ".embedded-postgres-go/extracted"))
	ui.Ok()
}

func stopPostgres(binaryExtractLocation string) {
	postgresBinary := filepath.Join(binaryExtractLocation, "bin/pg_ctl")
	dataLocation := filepath.Join(binaryExtractLocation, "data")
	postgresProcess := exec.Command(postgresBinary, "stop", "-w",
		"-D", dataLocation)
	postgresProcess.Run()
}
