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
		Short: "Create a boundary with its transformations",
		Long: `
A boundary allows to map code folders with tags. 

You can use it to categorize an application. For example, you can define an architectural boundary with
the different layers of an application. Or you can define a boundary for production code vs test code.
`,
		Example: "gocan create-boundary myboundary --scene myscene --app myapp --transformation src:src/main/ --transformation test:src/test/",
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

			a, err := app.FindAppByAppNameAndSceneName(connection, appName, sceneName)
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
		Short: "List the boundaries defined for an application",
		Example: "gocan boundaries --app myapp --scene myscene",
		Args: cobra.NoArgs,
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

			a, err := app.FindAppByAppNameAndSceneName(connection, appName, sceneName)
			if err != nil {
				return errors.Wrap(err, "Application not found")
			}

			ctx.Ui.Say("Retrieving boundaries...")

			data, err := c.QueryByAppId(a.Id)
			if err != nil {
				return errors.Wrap(err, "Unable to fetch the boundaries")
			}

			ctx.Ui.Ok()

			if len(data) > 0 {
				table := ctx.Ui.Table([]string{"id", "name", "transformations"})
				for _, b := range data {
					transformations := ""
					for _, t := range b.Transformations {
						transformations += t.Name + ":" + t.Path + " | "
					}
					table.Add(b.Id, b.Name, transformations)
				}
				table.Print()
			} else {
				ctx.Ui.Say("No boundaries found.")
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&sceneName, "scene", "s", "", "Scene name")
	cmd.Flags().StringVarP(&appName, "app", "a", "", "App name")

	return cmd
}