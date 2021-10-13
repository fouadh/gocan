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

	commands := []*cobra.Command{
		web_ui.NewStartUiCommand(ctx),
		db2.NewSetupDbCommand(ctx),
		db2.NewStartDbCommand(ctx),
		db2.NewStopDbCommand(ctx),
		db2.NewMigrateDb(ctx),
		scene.NewCreateScene(ctx),
		scene.NewDeleteScene(ctx),
		scene.NewScenes(ctx),
		app.NewCreateAppCommand(ctx),
		app.NewAppsCommand(ctx),
		app.NewAppSummary(ctx),
		app.NewDeleteApp(ctx),
		history.NewImportHistoryCommand(ctx),
		revision.NewRevisionsCommand(ctx),
		revision.NewHotspotsCommand(ctx),
		revision.NewRevisionTrends(ctx),
		coupling.NewCouplingCommand(ctx),
		coupling.NewSocCommand(ctx),
		coupling.NewCouplingHierarchyCommand(ctx),
		developer.NewMainDevelopers(ctx),
		developer.NewEntityEfforts(ctx),
		developer.NewKnowledgeMapCommand(ctx),
		developer.NewDevsCommand(ctx),
		churn.NewCodeChurn(ctx),
		modus_operandi.NewModusOperandi(ctx),
		active_set.NewActiveSet(ctx),
		boundary.NewCreateBoundary(ctx),
		boundary.NewDeleteBoundary(ctx),
		boundary.NewBoundaries(ctx),
		complexity.NewCreateComplexityAnalysis(ctx),
		complexity.NewDeleteComplexityAnalysis(ctx),
		complexity.NewComplexityAnalyses(ctx),
	}

	rootCmd.AddCommand(commands...)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
