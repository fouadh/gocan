package main

import (
	create_scene "com.fha.gocan/internal/create-scene"
	"com.fha.gocan/internal/platform/db"
	"com.fha.gocan/internal/platform/terminal"
	setup_db "com.fha.gocan/internal/setup-db"
	start_db "com.fha.gocan/internal/start-db"
	stop_db "com.fha.gocan/internal/stop-db"
	"com.fha.gocan/internal/ui"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "gocan",
}

const dsn = "host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable"

var uiCmd = ui.BuildUiCommand()

func main() {
	ui := terminal.NewUI(rootCmd.OutOrStdout(), rootCmd.ErrOrStderr())
	dataSource := db.SqlxDataSource{
		Dsn: dsn,
	}

	var createCmd = create_scene.BuildCreateSceneCmd(&dataSource, ui)
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(uiCmd)
	rootCmd.AddCommand(setup_db.BuildSetupDbCmd(ui))
	rootCmd.AddCommand(start_db.BuildStartDbCmd(ui))
	rootCmd.AddCommand(stop_db.BuildStopDbCmd(ui))

	rootCmd.Execute()
}

func connect(dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", dsn)
	return db, err
}

