package main

import (
	"com.fha.gocan/business/core/app"
	"com.fha.gocan/business/core/churn"
	"com.fha.gocan/business/core/coupling"
	db2 "com.fha.gocan/business/core/db"
	"com.fha.gocan/business/core/developer"
	"com.fha.gocan/business/core/history"
	"com.fha.gocan/business/core/revision"
	"com.fha.gocan/business/core/scene"
	web_ui "com.fha.gocan/business/core/ui"
	context "com.fha.gocan/foundation"
	"com.fha.gocan/foundation/db"
	"com.fha.gocan/foundation/terminal"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use: "gocan",
}

func main() {
	ui := terminal.NewUI(rootCmd.OutOrStdout(), rootCmd.ErrOrStderr())
	config, err := db.ReadConfig()
	if err  != nil {
		ui.Failed(errors.Cause(err).Error())
		os.Exit(2)
	}
	ctx := context.New(ui, config)

	rootCmd.AddCommand(web_ui.NewStartUiCommand(ctx))
	rootCmd.AddCommand(db2.NewSetupDbCommand(ctx))
	rootCmd.AddCommand(db2.NewStartDbCommand(ctx))
	rootCmd.AddCommand(db2.NewStopDbCommand(ctx))
	rootCmd.AddCommand(scene.NewCreateSceneCommand(ctx))
	rootCmd.AddCommand(scene.NewScenesCommand(ctx))
	rootCmd.AddCommand(app.NewCreateAppCommand(ctx))
	rootCmd.AddCommand(app.NewAppsCommand(ctx))
	rootCmd.AddCommand(history.NewImportHistoryCommand(ctx))
	rootCmd.AddCommand(revision.NewRevisionsCommand(*ctx))
	rootCmd.AddCommand(revision.NewHotspotsCommand(*ctx))
	rootCmd.AddCommand(coupling.NewCouplingCommand(ctx))
	rootCmd.AddCommand(coupling.NewSocCommand(*ctx))
	rootCmd.AddCommand(developer.NewMainDevelopers(*ctx))
	rootCmd.AddCommand(developer.NewKnowledgeMapCommand(*ctx))
	rootCmd.AddCommand(churn.NewCodeChurn(*ctx))

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}