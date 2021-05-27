package main

import (
	create_scene "com.fha.gocan/internal/create-scene"
	context "com.fha.gocan/internal/platform"
	"com.fha.gocan/internal/platform/terminal"
	setup_db "com.fha.gocan/internal/setup-db"
	start_db "com.fha.gocan/internal/start-db"
	stop_db "com.fha.gocan/internal/stop-db"
	"com.fha.gocan/internal/ui"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "gocan",
}

const dsn = "host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable"

var uiCmd = ui.BuildUiCommand()

func main() {
	ui := terminal.NewUI(rootCmd.OutOrStdout(), rootCmd.ErrOrStderr())
	ctx := context.New(dsn, ui)

	var createCmd = create_scene.BuildCreateSceneCmd(ctx)
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(uiCmd)
	rootCmd.AddCommand(setup_db.BuildSetupDbCmd(ctx))
	rootCmd.AddCommand(start_db.BuildStartDbCmd(ctx))
	rootCmd.AddCommand(stop_db.BuildStopDbCmd(ctx))

	rootCmd.Execute()
}