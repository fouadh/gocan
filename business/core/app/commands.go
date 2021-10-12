package app

import (
	"com.fha.gocan/business/core/commit"
	"com.fha.gocan/business/data/store/app"
	"com.fha.gocan/foundation"
	"com.fha.gocan/foundation/date"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"strconv"
)

func NewCreateAppCommand(ctx *foundation.Context) *cobra.Command {
	var sceneName string

	cmd := cobra.Command{
		Use:  "create-app",
		Short: "Create an application in a scene",
		Args: cobra.ExactArgs(1),
		Example: "gocan create-app myapp --scene myscene",
		RunE: func(cmd *cobra.Command, args []string) error {
			connection, err := ctx.GetConnection()
			if err != nil {
				return err
			}
			defer connection.Close()

			ctx.Ui.Log("Creating the app...")
			core := NewCore(connection)
			a, err := core.Create(*ctx, args[0], sceneName)

			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("Unable to create the app: %s", err.Error()))
			}

			ctx.Ui.Log(fmt.Sprintln("App", a.Id, "created."))
			ctx.Ui.Ok()
			return nil
		},
	}
	cmd.Flags().StringVarP(&sceneName, "scene", "s", "", "Scene name")
	return &cmd
}

func NewAppsCommand(ctx *foundation.Context) *cobra.Command {
	var sceneName string
	var csv bool

	cmd := cobra.Command{
		Use: "apps",
		Short: "List the applications associated with a scene",
		Args: cobra.NoArgs,
		Example: "gocan apps --scene myscene",
		RunE: func(cmd *cobra.Command, args []string) error {
			connection, err := ctx.GetConnection()
			if err != nil {
				return err
			}
			defer connection.Close()

			ctx.Ui.Log("Retrieving the apps...")
			core := NewCore(connection)
			apps, err := core.QueryBySceneName(sceneName)

			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("Unable to fetch the apps: %s", err.Error()))
			}

			ctx.Ui.Ok()

			if len(apps) > 0 {
				printApps(ctx, apps, csv)
			} else {
				ctx.Ui.Log("There is no application in this scene.")
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&sceneName, "scene", "s", "", "Scene name")
	cmd.Flags().BoolVar(&csv, "csv", false, "get the results in csv format")
	return &cmd
}

func NewDeleteApp(ctx foundation.Context) *cobra.Command {
	var sceneName string

	cmd := cobra.Command{
		Use: "delete-app",
		Aliases: []string{"da"},
		Short: "Delete an application",
		Example: "gocan delete-app myapp -s myscene",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			connection, err := ctx.GetConnection()
			if err != nil {
				return err
			}
			defer connection.Close()

			c := NewCore(connection)

			a, err := c.FindAppByAppNameAndSceneName(args[0], sceneName)
			if err != nil {
				return errors.Wrap(err, "Unable to retrieve the app")
			}

			ctx.Ui.Log("Deleting the app...")
			if err := c.Delete(a.Id); err != nil {
				return errors.Wrap(err, "Unable to delete the app")
			}
			ctx.Ui.Log("OK")

			return nil
		},
	}

	cmd.Flags().StringVarP(&sceneName, "scene", "s", "", "Scene name")

	return &cmd
}

func NewAppSummary(ctx foundation.Context) *cobra.Command {
	var sceneName string
	var before string
	var after string
	var csv bool

	cmd := cobra.Command{
		Use: "app-summary",
		Short: "Get an application summary information",
		Example: `
gocan app-summary myapp --scene myscene --after 2021-01-01 --before 2021-02-01
gocan app-summary myapp --scene myscene --before 2021-02-01
gocan app-summary myapp --scene myscene --after 2021-01-01
gocan app-summary myapp --scene myscene
`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			connection, err := ctx.GetConnection()
			if err != nil {
				return err
			}
			defer connection.Close()

			ctx.Ui.Log("Retrieving the apps...")
			c := NewCore(connection)
			ct := commit.NewCore(connection)

			a, err := c.FindAppByAppNameAndSceneName(args[0], sceneName)
			if err != nil {
				return errors.Wrap(err, "Invalid app")
			}

			beforeTime, afterTime, err := ct.ExtractDateRangeFromArgs(a.Id, before, after)
			if err != nil {
				return errors.Wrap(err, "Invalid date range")
			}

			summary, err := c.QuerySummary(a.Id, beforeTime, afterTime)
			if err != nil {
				return errors.Wrap(err, "No summary found for app " + a.Id + " between " + date.FormatDay(afterTime) + " and " + date.FormatDay(beforeTime))
			}

			ctx.Ui.Ok()

			table := ctx.Ui.Table([]string{"id", "name", "commits", "entities", "entities-changed", "authors"}, csv)
			table.Add(summary.Id, summary.Name, strconv.Itoa(summary.NumberOfCommits), strconv.Itoa(summary.NumberOfEntities), strconv.Itoa(summary.NumberOfEntitiesChanged), strconv.Itoa(summary.NumberOfAuthors))
			table.Print()

			return nil
		},
	}

	cmd.Flags().StringVarP(&sceneName, "scene", "s", "", "Scene name")
	cmd.Flags().StringVarP(&before, "before", "", "", "Fetch the summary of coupling before this day")
	cmd.Flags().StringVarP(&after, "after", "", "", "Fetch all the summary of coupling after this day")
	cmd.Flags().BoolVar(&csv, "csv", false, "get the results in csv format")
	return &cmd
}

func printApps(ctx *foundation.Context, apps []app.App, csv bool) {
	table := ctx.Ui.Table([]string{
		"id",
		"name",
	}, csv)

	for _, a := range apps {
		table.Add(a.Id, a.Name)
	}

	table.Print()
}
