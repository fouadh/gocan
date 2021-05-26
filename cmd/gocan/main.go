package main

import (
	create_scene "com.fha.gocan/internal/create-scene"
	"com.fha.gocan/internal/db"
	"com.fha.gocan/internal/terminal"
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
	database := db.EmbeddedDatabase{}
	database.Start(ui)
	defer database.Stop(ui)
	db.Migrate(dsn, ui)

	connection := database.Connect(dsn, ui)
	var createCmd = create_scene.BuildCreateSceneCmd(connection, ui)

	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(uiCmd)

	rootCmd.Execute()
}

func connect(dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", dsn)
	return db, err
}

