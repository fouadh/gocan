package boundary

import (
	"com.fha.gocan/business/core/app"
	"com.fha.gocan/foundation"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func Commands(ctx foundation.Context) []*cobra.Command {
	return []*cobra.Command{
		create(ctx),
		delete(ctx),
		list(ctx),
	}
}

func create(ctx foundation.Context) *cobra.Command {
	var sceneName string
	var appName string
	var modules []string
	var verbose bool

	cmd := &cobra.Command{
		Use:   "create-boundaries",
		Short: "Create a boundary with its modules",
		Long: `
A boundary allows to map code folders with tags. 

You can use it to categorize an application. For example, you can define an architectural boundary with
the different layers of an application. Or you can define a boundary for production code vs test code.
`,
		Example: "gocan create-boundaries myboundary --scene myscene --app myapp --module src:src/main/ --module test:src/test/",
		Aliases: []string{"cb", "create-boundary"},
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if appName == "" {
				return fmt.Errorf("No application provided")
			}
			if sceneName == "" {
				return fmt.Errorf("No scene provided")
			}
			if len(modules) == 0 {
				return fmt.Errorf("No module provided")
			}

			ctx.Ui.SetVerbose(verbose)
			connection, err := ctx.GetConnection()
			if err != nil {
				return err
			}
			defer connection.Close()

			c := NewCore(connection)

			a, err := app.FindAppByAppNameAndSceneName(connection, appName, sceneName)
			if err != nil {
				return errors.Wrap(err, "Application not found")
			}

			ctx.Ui.Log("Creating boundary...")
			_, err = c.Create(a.Id, args[0], modules)
			if err != nil {
				return err
			}
			ctx.Ui.Print("Boundary has been created.")
			ctx.Ui.Ok()

			return nil
		},
	}

	cmd.Flags().StringVarP(&sceneName, "scene", "s", "", "Scene name")
	cmd.Flags().StringVarP(&appName, "app", "a", "", "App name")
	cmd.Flags().StringSliceVarP(&modules, "module", "t", nil, "Modules")
	cmd.Flags().BoolVar(&verbose, "verbose", false, "display the log information")

	cmd.MarkFlagRequired("scene")
	cmd.MarkFlagRequired("app")
	cmd.MarkFlagRequired("module")

	return cmd
}

func delete(ctx foundation.Context) *cobra.Command {
	var sceneName string
	var appName string
	var verbose bool

	cmd := &cobra.Command{
		Use:     "delete-boundary",
		Aliases: []string{"db"},
		Short:   "Delete an application boundary",
		Example: "gocan delete-boundary myboundary --app myapp --scene myscene",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if appName == "" {
				return fmt.Errorf("No application provided")
			}
			if sceneName == "" {
				return fmt.Errorf("No scene provided")
			}

			ctx.Ui.SetVerbose(verbose)
			connection, err := ctx.GetConnection()
			if err != nil {
				return err
			}
			defer connection.Close()

			c := NewCore(connection)

			a, err := app.FindAppByAppNameAndSceneName(connection, appName, sceneName)
			if err != nil {
				return errors.Wrap(err, "Application not found")
			}

			ctx.Ui.Log("Deleting boundary...")

			if err := c.DeleteBoundaryByName(a.Id, args[0]); err != nil {
				return errors.Wrap(err, "Unable to delete boundary")
			}

			ctx.Ui.Print("Boundary has been deleted.")
			ctx.Ui.Ok()
			return nil
		},
	}

	cmd.Flags().StringVarP(&sceneName, "scene", "s", "", "Scene name")
	cmd.Flags().StringVarP(&appName, "app", "a", "", "App name")
	cmd.Flags().BoolVar(&verbose, "verbose", false, "display the log information")

	cmd.MarkFlagRequired("scene")
	cmd.MarkFlagRequired("app")

	return cmd
}

func list(ctx foundation.Context) *cobra.Command {
	var sceneName string
	var appName string
	var verbose bool

	cmd := &cobra.Command{
		Use:     "boundaries",
		Short:   "List the boundaries defined for an application",
		Example: "gocan boundaries --app myapp --scene myscene",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if appName == "" {
				return fmt.Errorf("No application provided")
			}
			if sceneName == "" {
				return fmt.Errorf("No scene provided")
			}

			ctx.Ui.SetVerbose(verbose)
			connection, err := ctx.GetConnection()
			if err != nil {
				return err
			}
			defer connection.Close()

			c := NewCore(connection)

			a, err := app.FindAppByAppNameAndSceneName(connection, appName, sceneName)
			if err != nil {
				return errors.Wrap(err, "Application not found")
			}

			ctx.Ui.Log("Retrieving boundaries...")

			data, err := c.QueryByAppId(a.Id)
			if err != nil {
				return errors.Wrap(err, "Unable to fetch the boundaries")
			}

			ctx.Ui.Ok()

			if len(data) > 0 {
				table := ctx.Ui.Table([]string{"id", "name", "modules"}, false)
				for _, b := range data {
					modules := ""
					for _, m := range b.Modules {
						modules += m.Name + ":" + m.Path + " | "
					}
					table.Add(b.Id, b.Name, modules)
				}
				table.Print()
			} else {
				ctx.Ui.Log("No boundaries found.")
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&sceneName, "scene", "s", "", "Scene name")
	cmd.Flags().StringVarP(&appName, "app", "a", "", "App name")
	cmd.Flags().BoolVar(&verbose, "verbose", false, "display the log information")

	cmd.MarkFlagRequired("scene")
	cmd.MarkFlagRequired("app")

	return cmd
}
