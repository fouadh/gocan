package db

import (
	"com.fha.gocan/business/platform/config"
	"com.fha.gocan/business/platform/terminal"
	"fmt"
	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"log"
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
		Port(uint32(ed.Config.Port)).
		DataPath(ed.Config.EmbeddedDataPath))

	ui.Say("Starting the embedded database...")
	if err := ed.database.Start(); err != nil {
		ui.Failed(fmt.Sprintf("Cannot start the database: %+v\n", err))
		os.Exit(1)
	}
	ui.Ok()
	ui.Say("Applying migrations...")
	Migrate(ed.Config.Dsn(), ui)
	ui.Ok()
}

func (ed *EmbeddedDatabase) Stop(ui terminal.UI) {
	ui.Say("Stopping the database...")
	usr, _ := user.Current()
	dir := usr.HomeDir
	// todo the path to the db should come from the configuration
	stopPostgres(filepath.Join(dir, ".embedded-postgres-go/extracted"), ed.Config)
	ui.Ok()
}

func startPostgres(binaryExtractLocation string, c *config.Config) error {
	postgresBinary := filepath.Join(binaryExtractLocation, "bin/pg_ctl")
	postgresProcess := exec.Command(postgresBinary, "start", "-w",
		"-D", c.EmbeddedDataPath,
		"-o", fmt.Sprintf(`"-p %d"`, c.Port))
	log.Println(postgresProcess.String())
	//postgresProcess.Stderr = config.logger
	//postgresProcess.Stdout = config.logger

	if err := postgresProcess.Run(); err != nil {
		return fmt.Errorf("could not start postgres using %s", postgresProcess.String())
	}

	return nil
}


func stopPostgres(binaryExtractLocation string, c *config.Config) {
	postgresBinary := filepath.Join(binaryExtractLocation, "bin/pg_ctl")
	postgresProcess := exec.Command(postgresBinary, "stop", "-w",
		"-D", c.EmbeddedDataPath)
	postgresProcess.Run()
}
