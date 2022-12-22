package main

import (
	"com.fha.gocan/business/core/age"
	"com.fha.gocan/business/core/app"
	"com.fha.gocan/business/core/boundary"
	"com.fha.gocan/business/core/churn"
	"com.fha.gocan/business/core/complexity"
	"com.fha.gocan/business/core/configuration"
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
	"github.com/spf13/cobra/doc"
	"os"
)

var Version = "development"

func main() {
	rootCmd := Root()
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func Root() *cobra.Command {
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
		age.Commands(ctx),
		developer.Commands(ctx),
		boundary.Commands(ctx),
		complexity.Commands(ctx),
		churn.Commands(ctx),
		modus_operandi.Commands(ctx),
		storyboard.Commands(ctx),
		configuration.Commands(ctx),
	}

	for _, c := range commands {
		rootCmd.AddCommand(c...)
	}
	rootCmd.AddCommand(GenerateDoc(), Completion())
	return rootCmd
}

func Completion() *cobra.Command {
	return &cobra.Command{
		Use:       "completion",
		Short:     "Generate completion script",
		ValidArgs: []string{"bash", "zsh", "fish", "powershell"},
		RunE: func(cmd *cobra.Command, args []string) error {
			switch args[0] {
			case "bash":
				return cmd.Root().GenBashCompletion(os.Stdout)
			case "zsh":
				return cmd.Root().GenZshCompletion(os.Stdout)
			case "fish":
				return cmd.Root().GenFishCompletion(os.Stdout, true)
			case "powershell":
				return cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
			}
			return nil
		},
	}
}

func GenerateDoc() *cobra.Command {
	var directory string

	cmd := cobra.Command{
		Use: "generate-doc",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := doc.GenMarkdownTree(cmd.Root(), directory); err != nil {
				return errors.Wrap(err, "Unable to generate the doc")
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&directory, "directory", "d", ".", "Directory where the doc will be generated")
	return &cmd
}
