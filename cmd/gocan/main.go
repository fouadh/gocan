package main

import (
	create_app "com.fha.gocan/internal/create-app"
	create_scene "com.fha.gocan/internal/create-scene"
	"com.fha.gocan/internal/import_history"
	context "com.fha.gocan/internal/platform"
	"com.fha.gocan/internal/platform/config"
	"com.fha.gocan/internal/platform/terminal"
	setup_db "com.fha.gocan/internal/setup-db"
	start_db "com.fha.gocan/internal/start-db"
	stop_db "com.fha.gocan/internal/stop-db"
	ui2 "com.fha.gocan/internal/ui"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use: "gocan",
}

func main() {
	ui := terminal.NewUI(rootCmd.OutOrStdout(), rootCmd.ErrOrStderr())
	config, err := config.ReadConfig()
	if err  != nil {
		ui.Failed(errors.Cause(err).Error())
		os.Exit(2)
	}
	ctx := context.New(ui, config)

	rootCmd.AddCommand(ui2.NewCommand(ctx))
	rootCmd.AddCommand(setup_db.NewCommand(ctx))
	rootCmd.AddCommand(start_db.NewCommand(ctx))
	rootCmd.AddCommand(stop_db.NewCommand(ctx))
	rootCmd.AddCommand(create_scene.NewCommand(ctx))
	rootCmd.AddCommand(create_app.NewCommand(ctx))
	rootCmd.AddCommand(import_history.NewCommand(ctx))

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}