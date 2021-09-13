package app

import (
	"com.fha.gocan/business/data/store/app"
	"com.fha.gocan/foundation"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func NewCreateAppCommand(ctx *foundation.Context) *cobra.Command {
	var sceneName string

	cmd := cobra.Command{
		Use:  "create-app",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			connection, err := ctx.GetConnection()
			if err != nil {
				return err
			}

			ctx.Ui.Say("Creating the app...")
			core := NewCore(connection)
			a, err := core.Create(*ctx, args[0], sceneName)

			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("Unable to create the app: %s", err.Error()))
			}

			ctx.Ui.Say(fmt.Sprintln("App", a.Id, "created."))
			ctx.Ui.Ok()
			return nil
		},
	}
	cmd.Flags().StringVarP(&sceneName, "scene", "s", "", "Scene name")
	return &cmd
}

func NewAppsCommand(ctx *foundation.Context) *cobra.Command {
	var sceneName string

	cmd := cobra.Command{
		Use: "apps",
		RunE: func(cmd *cobra.Command, args []string) error {
			connection, err := ctx.GetConnection()
			if err != nil {
				return err
			}

			ctx.Ui.Say("Retrieving the apps...")
			core := NewCore(connection)
			apps, err := core.QueryBySceneName(sceneName)

			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("Unable to fetch the apps: %s", err.Error()))
			}

			ctx.Ui.Ok()

			printApps(ctx, apps)

			return nil
		},
	}

	cmd.Flags().StringVarP(&sceneName, "scene", "s", "", "Scene name")
	return &cmd
}

func printApps(ctx *foundation.Context, apps []app.App) {
	table := ctx.Ui.Table([]string{
		"id",
		"name",
	})

	for _, a := range apps {
		table.Add(a.Id, a.Name)
	}

	table.Print()
}