package revision

import (
	"com.fha.gocan/business/core"
	"com.fha.gocan/business/core/app"
	"com.fha.gocan/business/data/store/revision"
	context "com.fha.gocan/foundation"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"sort"
	"strconv"
)

func Commands(ctx context.Context) []*cobra.Command {
	return []*cobra.Command{
		list(ctx),
		authors(ctx),
		createTrends(ctx),
		trends(ctx),
	}
}

func authors(ctx context.Context) *cobra.Command {
	var sceneName string
	var before string
	var after string
	var boundaryName string
	var csv bool
	var verbose bool

	cmd := cobra.Command{
		Use:   "revisions-authors",
		Args:  cobra.ExactArgs(1),
		Short: "Get the entities of an application ordered by their number of authors",
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

			var revisions []revision.Revision

			if boundaryName == "" {
				revisions, err = c.Query(a.Id, beforeTime, afterTime)
			} else {
				b, err := c.boundary.QueryByAppIdAndName(a.Id, boundaryName)
				if err != nil {
					return errors.Wrap(err, "Unable to get boundary")
				}
				revisions, err = c.QueryByBoundary(a.Id, b, beforeTime, afterTime)
			}

			if err != nil {
				ui.Failed("Cannot fetch revisions: " + err.Error())
				return err
			}

			ui.Ok()

			sort.Slice(revisions, func(i, j int) bool {
				return revisions[i].NumberOfAuthors > revisions[j].NumberOfAuthors
			})

			table := ui.Table([]string{
				"entity",
				"n-authors",
				"n-revs",
			}, csv)
			for _, revision := range revisions {
				table.Add(revision.Entity, strconv.Itoa(revision.NumberOfAuthors), strconv.Itoa(revision.NumberOfRevisions))
			}
			table.Print()
			return nil
		},
	}

	cmd.Flags().StringVarP(&sceneName, "scene", "s", "", "Scene name")
	cmd.Flags().StringVarP(&boundaryName, "boundary", "", "", "Optional boundary name to get the analysis for a specific boundary")
	cmd.Flags().StringVarP(&before, "before", "", "", "Fetch all the revisions before this day")
	cmd.Flags().StringVarP(&after, "after", "", "", "Fetch all the revisions after this day")
	cmd.Flags().BoolVar(&csv, "csv", false, "get the results in csv format")
	cmd.Flags().BoolVar(&verbose, "verbose", false, "display the log information")

	cmd.MarkFlagRequired("scene")

	return &cmd
}

func list(ctx context.Context) *cobra.Command {
	var sceneName string
	var before string
	var after string
	var boundaryName string
	var moduleName string
	var csv bool
	var verbose bool

	cmd := cobra.Command{
		Use:   "revisions",
		Args:  cobra.ExactArgs(1),
		Short: "Get the entities of an application ordered by their number of revisions",
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

			var revisions []revision.Revision

			if boundaryName == "" {
				revisions, err = c.Query(a.Id, beforeTime, afterTime)
			} else if moduleName != "" {
				b, err := c.boundary.QueryByAppIdAndName(a.Id, boundaryName)
				if err != nil {
					return errors.Wrap(err, "Unable to get boundary")
				}

				var mod = b.FindModule(moduleName)
				if mod.Name == "" {
					return errors.New("unable to retrieve module " + moduleName)
				}

				revisions, err = c.QueryByModule(a.Id, mod, beforeTime, afterTime)
			} else {
				b, err := c.boundary.QueryByAppIdAndName(a.Id, boundaryName)
				if err != nil {
					return errors.Wrap(err, "Unable to get boundary")
				}
				revisions, err = c.QueryByBoundary(a.Id, b, beforeTime, afterTime)
			}

			if err != nil {
				ui.Failed("Cannot fetch revisions: " + err.Error())
				return err
			}

			ui.Ok()

			table := ui.Table([]string{
				"entity",
				"n-revs",
				"code",
			}, csv)
			for _, revision := range revisions {
				table.Add(revision.Entity, strconv.Itoa(revision.NumberOfRevisions), strconv.Itoa(revision.Code))
			}
			table.Print()
			return nil
		},
	}

	cmd.Flags().StringVarP(&sceneName, "scene", "s", "", "Scene name")
	cmd.Flags().StringVarP(&boundaryName, "boundary", "", "", "Optional boundary name to get the analysis for a specific boundary")
	cmd.Flags().StringVarP(&moduleName, "module", "", "", "Optional module name to get the analysis for a specific module. Boundary must also be set.")
	cmd.Flags().StringVarP(&before, "before", "", "", "Fetch all the revisions before this day")
	cmd.Flags().StringVarP(&after, "after", "", "", "Fetch all the revisions after this day")
	cmd.Flags().BoolVar(&csv, "csv", false, "get the results in csv format")
	cmd.Flags().BoolVar(&verbose, "verbose", false, "display the log information")

	cmd.MarkFlagRequired("scene")

	return &cmd
}

func createTrends(ctx context.Context) *cobra.Command {
	var sceneName string
	var appName string
	var boundaryName string
	var before string
	var after string
	var verbose bool

	cmd := cobra.Command{
		Use:     "create-revisions-trends",
		Aliases: []string{"crt"},
		RunE: func(cmd *cobra.Command, args []string) error {
			ui := ctx.Ui
			ui.SetVerbose(verbose)

			connection, err := ctx.GetConnection()
			if err != nil {
				return err
			}
			defer connection.Close()

			c := NewCore(connection)
			a, beforeTime, afterTime, err := core.ExtractDateRangeAndAppFromArgs(connection, sceneName, appName, before, after)
			if err != nil {
				return errors.Wrap(err, "Command failed")
			}

			b, err := c.boundary.QueryByAppIdAndName(a.Id, boundaryName)
			if err != nil {
				return errors.Wrap(err, "Boundary not found")
			}

			ui.Log("Creating revisions trends...")
			if err := c.CreateRevisionTrends(args[0], a.Id, b, beforeTime, afterTime); err != nil {
				return errors.Wrap(err, "Cannot create revisions trends")
			}

			ui.Ok()

			return nil
		},
	}

	cmd.Flags().StringVarP(&sceneName, "scene", "s", "", "Scene name")
	cmd.Flags().StringVarP(&appName, "app", "a", "", "App name")
	cmd.Flags().StringVarP(&before, "before", "", "", "Create trends for before this day")
	cmd.Flags().StringVarP(&after, "after", "", "", "Create trends for after this day")
	cmd.Flags().StringVarP(&boundaryName, "boundary", "", "", "Boundary to use")
	cmd.Flags().BoolVar(&verbose, "verbose", false, "display the log information")

	cmd.MarkFlagRequired("scene")
	cmd.MarkFlagRequired("app")
	cmd.MarkFlagRequired("boundary")

	return &cmd
}

func trends(ctx context.Context) *cobra.Command {
	var sceneName string
	var appName string
	var csv bool
	var verbose bool

	cmd := cobra.Command{
		Use:     "revision-trends",
		Aliases: []string{"revisions-trends", "rt"},
		Short:   "Retrieve a saved revision trends",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ui := ctx.Ui
			ui.SetVerbose(verbose)

			connection, err := ctx.GetConnection()
			if err != nil {
				return err
			}
			defer connection.Close()

			c := NewCore(connection)
			a, err := app.FindAppByAppNameAndSceneName(connection, appName, sceneName)
			if err != nil {
				return errors.Wrap(err, "Command failed")
			}

			ui.Log("Getting revisions trends...")
			trends, err := c.RevisionTrendsByName(args[0], a.Id)
			if err != nil {
				return errors.Wrap(err, "Cannot get revisions trends")
			}

			ui.Ok()

			b, err := c.boundary.QueryById(trends.BoundaryId)
			if err != nil {
				return errors.Wrap(err, "Boundary not found")
			}
			headers := []string{"date"}
			for _, t := range b.Modules {
				headers = append(headers, t.Name)
			}
			table := ui.Table(headers, csv)
			for _, rt := range trends.Entries {
				cols := []string{rt.Date}
				for _, t := range b.Modules {
					cols = append(cols, strconv.Itoa(rt.FindEntityRevision(t.Name).NumberOfRevisions))
				}
				table.Add(cols...)
			}
			table.Print()

			return nil
		},
	}

	cmd.Flags().StringVarP(&sceneName, "scene", "s", "", "Scene name")
	cmd.Flags().StringVarP(&appName, "app", "a", "", "Application name")
	cmd.Flags().BoolVar(&csv, "csv", false, "get the results in csv format")
	cmd.Flags().BoolVar(&verbose, "verbose", false, "display the log information")

	cmd.MarkFlagRequired("scene")
	cmd.MarkFlagRequired("app")

	return &cmd
}
