package developer

import (
	"com.fha.gocan/business/core"
	"com.fha.gocan/business/core/app"
	"com.fha.gocan/foundation"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"strconv"
)

func Commands(ctx foundation.Context) []*cobra.Command {
	return []*cobra.Command{
		mainDevs(ctx),
		entityEffortsPerAuthor(ctx),
		entityEfforts(ctx),
		knowledgeMap(ctx),
		list(ctx),
		rename(ctx),
		createTeam(ctx),
	}
}

func createTeam(ctx foundation.Context) *cobra.Command {
	var sceneName string
	var appName string
	var members []string
	var verbose bool

	cmd := cobra.Command{
		Use:   "create-team",
		Short: "Create a team of developers with its members",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ui := ctx.Ui
			ui.SetVerbose(verbose)
			connection, err := ctx.GetConnection()
			if err != nil {
				return err
			}
			defer connection.Close()

			a, err := app.FindAppByAppNameAndSceneName(connection, appName, sceneName)
			if err != nil {
				return errors.Wrap(err, "error fetching the app")
			}
			if a.Id == "" {
				return errors.Errorf("Unable to retrieve the app")
			}

			c := NewCore(connection)

			teamName := args[0]
			if err := c.CreateTeam(a.Id, teamName, members); err != nil {
				return errors.Wrap(err, "Unable to create the team")
			}

			ui.Print("The team has been created.")

			ui.Ok()
			return nil
		},
	}

	cmd.Flags().StringVarP(&sceneName, "scene", "s", "", "Scene name")
	cmd.Flags().StringVarP(&appName, "app", "a", "", "App name")
	cmd.Flags().StringArrayVarP(&members, "member", "c", []string{}, "Team members")
	cmd.Flags().BoolVar(&verbose, "verbose", false, "display the log information")

	cmd.MarkFlagRequired("scene")
	cmd.MarkFlagRequired("app")
	cmd.MarkFlagRequired("members")

	return &cmd
}

func rename(ctx foundation.Context) *cobra.Command {
	var sceneName string
	var appName string
	var currentName string
	var newName string
	var verbose bool

	cmd := cobra.Command{
		Use:     "rename-dev",
		Aliases: []string{"rename-developer"},
		Short:   "Rename a developer in the database",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ui := ctx.Ui
			ui.SetVerbose(verbose)
			connection, err := ctx.GetConnection()
			if err != nil {
				return err
			}
			defer connection.Close()

			ui.Log("Renaming developer...")

			a, err := app.FindAppByAppNameAndSceneName(connection, appName, sceneName)
			if err != nil {
				return errors.Wrap(err, "error fetching the app")
			}
			if a.Id == "" {
				return errors.Errorf("Unable to retrieve the app")
			}

			c := NewCore(connection)

			if err := c.RenameDeveloper(a.Id, currentName, newName); err != nil {
				return errors.Wrap(err, "Unable to rename the developer")
			}

			ui.Print("The author has been renamed.")

			ui.Ok()
			return nil
		},
	}

	cmd.Flags().StringVarP(&sceneName, "scene", "s", "", "Scene name")
	cmd.Flags().StringVarP(&appName, "app", "a", "", "App name")
	cmd.Flags().StringVarP(&currentName, "current", "c", "", "Current developer name")
	cmd.Flags().StringVarP(&newName, "new", "n", "", "New developer name")
	cmd.Flags().BoolVar(&verbose, "verbose", false, "display the log information")

	cmd.MarkFlagRequired("scene")
	cmd.MarkFlagRequired("app")
	cmd.MarkFlagRequired("current")
	cmd.MarkFlagRequired("new")

	return &cmd
}

func mainDevs(ctx foundation.Context) *cobra.Command {
	var sceneName string
	var before string
	var after string
	var csv bool
	var verbose bool

	cmd := cobra.Command{
		Use:     "main-developers",
		Aliases: []string{"main_developers", "mainDevelopers", "md", "main-devs", "main_devs", "mainDevs"},
		Short:   "Get the main developers of an application entities",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ui := ctx.Ui
			ui.SetVerbose(verbose)
			connection, err := ctx.GetConnection()
			if err != nil {
				return err
			}
			defer connection.Close()

			ui.Log("Retrieving main developers...")

			c := NewCore(connection)

			a, beforeTime, afterTime, err := core.ExtractDateRangeAndAppFromArgs(connection, sceneName, args[0], before, after)
			if err != nil {
				return errors.Wrap(err, "Invalid argument(s)")
			}

			data, err := c.QueryMainDevelopers(a.Id, beforeTime, afterTime)
			if err != nil {
				return errors.Wrap(err, "Cannot retrieve main developers")
			}

			ui.Ok()

			table := ui.Table([]string{"entity", "main-dev", "added", "total-added", "ownership"}, csv)
			for _, dev := range data {
				table.Add(dev.Entity, dev.Author, fmt.Sprint(dev.Added), fmt.Sprint(dev.TotalAdded), fmt.Sprintf("%.2f", dev.Ownership))
			}
			table.Print()

			return nil
		},
	}

	cmd.Flags().StringVarP(&sceneName, "scene", "s", "", "Scene name")
	cmd.Flags().StringVarP(&before, "before", "", "", "Fetch the main developers before this day")
	cmd.Flags().StringVarP(&after, "after", "", "", "Fetch all the main developers after this day")
	cmd.Flags().BoolVar(&csv, "csv", false, "get the results in csv format")
	cmd.Flags().BoolVar(&verbose, "verbose", false, "display the log information")

	cmd.MarkFlagRequired("scene")

	return &cmd
}

func entityEffortsPerAuthor(ctx foundation.Context) *cobra.Command {
	var sceneName string
	var before string
	var after string
	var csv bool
	var verbose bool

	cmd := cobra.Command{
		Use:     "entity-efforts-per-author",
		Aliases: []string{"eepa"},
		Short:   "Get the efforts associated with entities of an application per author",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ui := ctx.Ui
			ui.SetVerbose(verbose)
			connection, err := ctx.GetConnection()
			if err != nil {
				return err
			}
			defer connection.Close()

			ui.Log("Retrieving entity efforts...")

			c := NewCore(connection)

			a, beforeTime, afterTime, err := core.ExtractDateRangeAndAppFromArgs(connection, sceneName, args[0], before, after)
			if err != nil {
				return errors.Wrap(err, "Invalid argument(s)")
			}

			data, err := c.QueryEntityEffortsPerAuthor(a.Id, beforeTime, afterTime)
			if err != nil {
				return errors.Wrap(err, "Cannot retrieve entity efforts per author")
			}

			ui.Ok()

			table := ui.Table([]string{"entity", "author", "author-revs", "total-revs"}, csv)
			for _, dev := range data {
				table.Add(dev.Entity, dev.Author, fmt.Sprint(dev.AuthorRevisions), fmt.Sprint(dev.TotalRevisions))
			}
			table.Print()

			return nil
		},
	}

	cmd.Flags().StringVarP(&sceneName, "scene", "s", "", "Scene name")
	cmd.Flags().StringVarP(&before, "before", "", "", "Fetch the entity efforts before this day")
	cmd.Flags().StringVarP(&after, "after", "", "", "Fetch all the entity efforts after this day")
	cmd.Flags().BoolVar(&csv, "csv", false, "get the results in csv format")
	cmd.Flags().BoolVar(&verbose, "verbose", false, "display the log information")

	cmd.MarkFlagRequired("scene")

	return &cmd
}
func entityEfforts(ctx foundation.Context) *cobra.Command {
	var sceneName string
	var before string
	var after string
	var csv bool
	var verbose bool

	cmd := cobra.Command{
		Use:     "entity-efforts",
		Aliases: []string{"entity_efforts", "entityEfforts", "ee"},
		Short:   "Get the cumulated efforts associated with entities of an application",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ui := ctx.Ui
			ui.SetVerbose(verbose)
			connection, err := ctx.GetConnection()
			if err != nil {
				return err
			}
			defer connection.Close()

			ui.Log("Retrieving entity efforts...")

			c := NewCore(connection)

			a, beforeTime, afterTime, err := core.ExtractDateRangeAndAppFromArgs(connection, sceneName, args[0], before, after)
			if err != nil {
				return errors.Wrap(err, "Invalid argument(s)")
			}

			ui.Log("Found app " + a.Id)

			data, err := c.QueryEntityEfforts(a.Id, beforeTime, afterTime)
			if err != nil {
				return errors.Wrap(err, "Cannot retrieve entity efforts")
			}

			ui.Ok()

			table := ui.Table([]string{"entity", "effort"}, csv)
			for _, dev := range data {
				table.Add(dev.Entity, fmt.Sprint(dev.Effort))
			}
			table.Print()

			return nil
		},
	}

	cmd.Flags().StringVarP(&sceneName, "scene", "s", "", "Scene name")
	cmd.Flags().StringVarP(&before, "before", "", "", "Fetch the entity efforts before this day")
	cmd.Flags().StringVarP(&after, "after", "", "", "Fetch all the entity efforts after this day")
	cmd.Flags().BoolVar(&csv, "csv", false, "get the results in csv format")
	cmd.Flags().BoolVar(&verbose, "verbose", false, "display the log information")

	cmd.MarkFlagRequired("scene")

	return &cmd
}

func knowledgeMap(ctx foundation.Context) *cobra.Command {
	var sceneName string
	var before string
	var after string
	var verbose bool

	cmd := cobra.Command{
		Use:     "knowledge-map",
		Aliases: []string{"knowledge_map", "knowledgeMap", "km"},
		Short:   "Get the knowledge map",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ui := ctx.Ui
			ui.SetVerbose(verbose)
			connection, err := ctx.GetConnection()
			if err != nil {
				return err
			}
			defer connection.Close()

			ui.Log("Building knowledge map...")

			c := NewCore(connection)

			a, beforeTime, afterTime, err := core.ExtractDateRangeAndAppFromArgs(connection, sceneName, args[0], before, after)
			if err != nil {
				return errors.Wrap(err, "Invalid argument(s)")
			}

			km, err := c.BuildKnowledgeMap(a, beforeTime, afterTime)

			ui.Ok()

			str, _ := json.MarshalIndent(km, "", "  ")
			ui.Print(string(str))

			return nil
		},
	}

	cmd.Flags().StringVarP(&sceneName, "scene", "s", "", "Scene name")
	cmd.Flags().StringVarP(&before, "before", "", "", "Calculate the knowledge map before this day")
	cmd.Flags().StringVarP(&after, "after", "", "", "Calculate all the knowledge map after this day")
	cmd.Flags().BoolVar(&verbose, "verbose", false, "display the log information")

	cmd.MarkFlagRequired("scene")

	return &cmd
}

func list(ctx foundation.Context) *cobra.Command {
	var sceneName string
	var before string
	var after string
	var csv bool
	var verbose bool

	cmd := cobra.Command{
		Use:     "devs",
		Aliases: []string{"developers"},
		Short:   "Get the developers of an application",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ui := ctx.Ui
			ui.SetVerbose(verbose)
			ui.Log("Getting app revisions...")

			connection, err := ctx.GetConnection()
			if err != nil {
				return err
			}
			defer connection.Close()

			c := NewCore(connection)
			a, beforeTime, afterTime, err := core.ExtractDateRangeAndAppFromArgs(connection, sceneName, args[0], before, after)

			devs, err := c.QueryDevelopers(a.Id, beforeTime, afterTime)

			if err != nil {
				return err
			}
			ui.Ok()

			table := ui.Table([]string{"name", "commits"}, csv)
			for _, dev := range devs {
				table.Add(dev.Name, strconv.Itoa(dev.NumberOfCommits))
			}
			table.Print()

			return nil
		},
	}

	cmd.Flags().StringVarP(&sceneName, "scene", "s", "", "Scene name")
	cmd.Flags().StringVarP(&before, "before", "", "", "Fetch the developers before this day")
	cmd.Flags().StringVarP(&after, "after", "", "", "Fetch all the developers after this day")
	cmd.Flags().BoolVar(&csv, "csv", false, "get the results in csv format")
	cmd.Flags().BoolVar(&verbose, "verbose", false, "display the log information")

	cmd.MarkFlagRequired("scene")

	return &cmd
}
