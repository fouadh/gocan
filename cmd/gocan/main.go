package main

import (
	create_scene "com.fha.gocan/internal/create-scene"
	context "com.fha.gocan/internal/platform"
	"com.fha.gocan/internal/platform/terminal"
	setup_db "com.fha.gocan/internal/setup-db"
	start_db "com.fha.gocan/internal/start-db"
	stop_db "com.fha.gocan/internal/stop-db"
	"com.fha.gocan/internal/ui"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use: "gocan",
}

var uiCmd = ui.BuildUiCommand()

func main() {
	ui := terminal.NewUI(rootCmd.OutOrStdout(), rootCmd.ErrOrStderr())
	config, err := setup_db.ReadConfig()
	if err  != nil {
		ui.Failed("Could not read the configuration file. Please use gocan setup-db to eventually regenerate it")
		os.Exit(2)
	}
	ctx := context.New(config.Dsn(), ui)

	var createCmd = create_scene.NewCommand(ctx)
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(uiCmd)
	rootCmd.AddCommand(setup_db.NewCommand(ctx))
	rootCmd.AddCommand(start_db.NewCommand(ctx))
	rootCmd.AddCommand(stop_db.NewCommand(ctx))

	if err := rootCmd.Execute(); err != nil {
		ui.Failed(errors.Cause(err).Error())
		os.Exit(1)
	}
}