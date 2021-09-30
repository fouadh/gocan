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
	Version: "0.1",
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
	rootCmd.AddCommand(scene.NewCreateScene(*ctx))
	rootCmd.AddCommand(scene.NewDeleteScene(*ctx))
	rootCmd.AddCommand(scene.NewScenes(*ctx))
	rootCmd.AddCommand(app.NewCreateAppCommand(ctx))
	rootCmd.AddCommand(app.NewAppsCommand(ctx))
	rootCmd.AddCommand(app.NewAppSummary(*ctx))
	rootCmd.AddCommand(app.NewDeleteApp(*ctx))
	rootCmd.AddCommand(history.NewImportHistoryCommand(ctx))
	rootCmd.AddCommand(revision.NewRevisionsCommand(*ctx))
	rootCmd.AddCommand(revision.NewHotspotsCommand(*ctx))
	rootCmd.AddCommand(revision.NewRevisionTrends(*ctx))
	rootCmd.AddCommand(coupling.NewCouplingCommand(ctx))
	rootCmd.AddCommand(coupling.NewSocCommand(*ctx))
	rootCmd.AddCommand(coupling.NewCouplingHierarchyCommand(*ctx))
	rootCmd.AddCommand(developer.NewMainDevelopers(*ctx))
	rootCmd.AddCommand(developer.NewEntityEfforts(*ctx))
	rootCmd.AddCommand(developer.NewKnowledgeMapCommand(*ctx))
	rootCmd.AddCommand(developer.NewDevsCommand(*ctx))
	rootCmd.AddCommand(churn.NewCodeChurn(*ctx))
	rootCmd.AddCommand(modus_operandi.NewModusOperandi(*ctx))
	rootCmd.AddCommand(active_set.NewActiveSet(*ctx))
	rootCmd.AddCommand(boundary.NewCreateBoundary(*ctx))
	rootCmd.AddCommand(boundary.NewDeleteBoundary(*ctx))
	rootCmd.AddCommand(boundary.NewBoundaries(*ctx))
	rootCmd.AddCommand(complexity.NewCreateComplexityAnalysis(*ctx))
	rootCmd.AddCommand(complexity.NewDeleteComplexityAnalysis(*ctx))
	rootCmd.AddCommand(complexity.NewComplexityAnalyses(*ctx))

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}