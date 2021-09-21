package boundary

import (
	"com.fha.gocan/business/core/app"
	"com.fha.gocan/foundation"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func NewCreateBoundary(ctx foundation.Context) *cobra.Command {
	var sceneName string
	var appName string
	var transformations []string

	cmd := &cobra.Command{
		Use:     "create-boundary",
		Aliases: []string{"cb"},
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if appName == "" {
				return fmt.Errorf("No application provided")
			}
			if sceneName == "" {
				return fmt.Errorf("No scene provided")
			}
			if len(transformations) == 0 {
				return fmt.Errorf("No transformation provided")
			}

			connection, err := ctx.GetConnection()
			if err != nil {
				return err
			}

			c := NewCore(connection)

			a, err := app.FindAppBySceneNameAndAppName(connection, appName, sceneName)
			if err != nil {
				return errors.Wrap(err, "Application not found")
			}

			ctx.Ui.Say("Creating boundary...")
			_, err = c.Create(a.Id, args[0], transformations)
			if err != nil {
				return err
			}
			ctx.Ui.Ok()

			return nil
		},
	}

	cmd.Flags().StringVarP(&sceneName, "scene", "s", "", "Scene name")
	cmd.Flags().StringVarP(&appName, "app", "a", "", "App name")
	cmd.Flags().StringSliceVarP(&transformations, "transformation", "t", nil, "Transformations")

	return cmd
}

func NewBoundaries(ctx foundation.Context) *cobra.Command {
	var sceneName string
	var appName string

	cmd := &cobra.Command{
		Use: "boundaries",
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			if appName == "" {
				return fmt.Errorf("No application provided")
			}
			if sceneName == "" {
				return fmt.Errorf("No scene provided")
			}

			connection, err := ctx.GetConnection()
			if err != nil {
				return err
			}

			c := NewCore(connection)

			a, err := app.FindAppBySceneNameAndAppName(connection, appName, sceneName)
			if err != nil {
				return errors.Wrap(err, "Application not found")
			}

			ctx.Ui.Say("Retrieving boundaries...")

			data, err := c.Query(a.Id)
			if err != nil {
				return errors.Wrap(err, "Unable to fetch the boundaries")
			}

			ctx.Ui.Ok()

			table := ctx.Ui.Table([]string{"id", "name", "transformations"})
			for _, b := range data {
				transformations := ""
				for _, t := range b.Transformations {
					transformations += t.Name + ":" + t.Path + " | "
				}
				table.Add(b.Id, b.Name, transformations)
			}
			table.Print()

			return nil
		},
	}

	cmd.Flags().StringVarP(&sceneName, "scene", "s", "", "Scene name")
	cmd.Flags().StringVarP(&appName, "app", "a", "", "App name")

	return cmd
}