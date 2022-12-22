package configuration

import (
	"com.fha.gocan/business/core/app"
	"com.fha.gocan/business/data/store/configuration"
	"com.fha.gocan/foundation"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"strings"
)

func Commands(ctx foundation.Context) []*cobra.Command {
	return []*cobra.Command{
		createExclusions(ctx),
		list(ctx),
	}
}

func list(ctx foundation.Context) *cobra.Command {
	var sceneName string
	var appName string
	var csv bool
	var verbose bool

	cmd := cobra.Command{
		Use:   "exclusions",
		Short: "Get exclusions associated with an app",
		RunE: func(cmd *cobra.Command, args []string) error {
			ui := ctx.Ui
			ui.SetVerbose(verbose)

			connection, err := ctx.GetConnection()
			if err != nil {
				return err
			}
			defer connection.Close()

			ctx.Ui.Log("Retrieving the application...")
			appCore := app.NewCore(connection)

			a, err := appCore.FindAppByAppNameAndSceneName(appName, sceneName)
			if err != nil {
				return errors.Wrap(err, "Invalid app")
			}

			ctx.Ui.Ok()

			configStore := configuration.NewStore(connection)

			ctx.Ui.Log("Retrieving exclusions...")

			exclusions, err := configStore.QueryExclusions(a.Id)

			if err != nil {
				return errors.Wrap(err, "Got an error while looking for exclusions!")
			}

			ctx.Ui.Ok()

			table := ui.Table([]string{"Exclusions"}, csv)
			for _, exclusion := range exclusions {
				table.Add(exclusion)
			}
			table.Print()

			return nil
		},
	}

	cmd.Flags().StringVarP(&sceneName, "scene", "s", "", "Scene name")
	cmd.Flags().StringVarP(&appName, "app", "a", "", "Application name")
	cmd.Flags().BoolVar(&csv, "csv", false, "get the results in csv format")
	cmd.Flags().BoolVar(&verbose, "verbose", false, "display the log information")

	return &cmd
}

func createExclusions(ctx foundation.Context) *cobra.Command {
	var sceneName string
	var appName string
	var verbose bool

	cmd := cobra.Command{
		Use:     "exclude",
		Short:   "Add files or folders to exclude from analysis",
		Example: "gocan exclude \"*.html;*.css;**/node-modules/**\" --app myapp --scene myscene",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ui := ctx.Ui
			ui.SetVerbose(verbose)

			exclusions := args[0]

			connection, err := ctx.GetConnection()
			if err != nil {
				return err
			}
			defer connection.Close()

			ctx.Ui.Log("Retrieving the application...")
			appCore := app.NewCore(connection)

			a, err := appCore.FindAppByAppNameAndSceneName(appName, sceneName)
			if err != nil {
				return errors.Wrap(err, "Invalid app")
			}

			ctx.Ui.Ok()

			configStore := configuration.NewStore(connection)

			ctx.Ui.Log("Adding exclusions to app...")
			if err := configStore.CreateExclusions(a.Id, strings.Split(exclusions, ";")); err != nil {
				return errors.Wrap(err, "Cannot create exclusions")
			}

			ctx.Ui.Ok()
			return nil
		},
	}

	cmd.Flags().StringVarP(&sceneName, "scene", "s", "", "Scene name")
	cmd.Flags().StringVarP(&appName, "app", "a", "", "Application name")
	cmd.Flags().BoolVar(&verbose, "verbose", false, "display the log information")

	return &cmd
}
