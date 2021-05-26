package main

import (
	create_scene "com.fha.gocan/internal/create-scene"
	init_db "com.fha.gocan/internal/init-db"
	"com.fha.gocan/internal/terminal"
	"com.fha.gocan/internal/ui"
	"fmt"
	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use: "gocan",
}

var uiCmd = ui.BuildUiCommand()
func main() {
	ui := terminal.NewUI(rootCmd.OutOrStdout(), rootCmd.ErrOrStderr())

	database := embeddedpostgres.NewDatabase()
	if err := database.Start(); err != nil {
		ui.Failed(fmt.Sprintf("Cannot start the database: %+v\n", err))
		os.Exit(1)
	}

	db, err := connect()
	if err != nil {
		ui.Failed(fmt.Sprintf("%s Cannot connect to the database: %+v\n", err))
		os.Exit(3)
	}

	defer func() {
		if err := database.Stop(); err != nil {
			ui.Failed(fmt.Sprintf("Cannot stop the database: %+v\n", err))
			os.Exit(2)
		}
	}()

	init_db.InitDb(ui)



	var createCmd = create_scene.BuildCreateSceneCmd(db, ui)
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(uiCmd)
	rootCmd.Execute()
}

func connect() (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", "host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable")
	return db, err
}

