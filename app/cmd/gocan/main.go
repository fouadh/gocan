package main

import (
	active_set "com.fha.gocan/business/core/active-set"
	"com.fha.gocan/business/core/app"
	"com.fha.gocan/business/core/boundary"
	"com.fha.gocan/business/core/churn"
	"com.fha.gocan/business/core/complexity"
	"com.fha.gocan/business/core/coupling"
	db2 "com.fha.gocan/business/core/db"
	"com.fha.gocan/business/core/developer"
	"com.fha.gocan/business/core/history"
	modus_operandi "com.fha.gocan/business/core/modus-operandi"
	"com.fha.gocan/business/core/revision"
	"com.fha.gocan/business/core/scene"
	"com.fha.gocan/business/core/storyboard"
	ui2 "com.fha.gocan/business/core/ui"
	context "com.fha.gocan/foundation"
	"com.fha.gocan/foundation/db"
	"com.fha.gocan/foundation/terminal"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"os"
)

var Version = "development"

func main() {
	rootCmd := &cobra.Command{
		Use:     "gocan",
		Version: Version,
	}

	ui := terminal.NewUI(rootCmd.OutOrStdout(), rootCmd.ErrOrStderr())
	config, err := db.ReadConfig()
	if err != nil {
		ui.Failed(errors.Cause(err).Error())
		os.Exit(2)
	}
	ctx := context.New(ui, config)

	commands := [][]*cobra.Command{
		ui2.Commands(ctx),
		db2.Commands(ctx),
		history.Commands(ctx),
		revision.Commands(ctx),
		coupling.Commands(ctx),
		scene.Commands(ctx),
		app.Commands(ctx),
		developer.Commands(ctx),
		boundary.Commands(ctx),
		complexity.Commands(ctx),
		churn.Commands(ctx),
		modus_operandi.Commands(ctx),
		active_set.Commands(ctx),
		storyboard.Commands(ctx),
	}

	for _, c := range commands {
		rootCmd.AddCommand(c...)
	}

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
