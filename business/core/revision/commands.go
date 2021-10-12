package revision

import (
	"com.fha.gocan/business/core"
	context "com.fha.gocan/foundation"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"strconv"
)

func NewRevisionsCommand(ctx context.Context) *cobra.Command {
	var sceneName string
	var before string
	var after string
	var csv bool
	var verbose bool

	cmd := cobra.Command{
		Use:  "revisions",
		Args: cobra.ExactArgs(1),
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

			revisions, err := c.Query(a.Id, beforeTime, afterTime)

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
	cmd.Flags().StringVarP(&before, "before", "a", "", "Fetch all the revisions before this day")
	cmd.Flags().StringVarP(&after, "after", "b", "", "Fetch all the revisions after this day")
	cmd.Flags().BoolVar(&csv, "csv", false, "get the results in csv format")
	cmd.Flags().BoolVar(&verbose, "verbose", false, "display the log information")
	return &cmd
}

func NewHotspotsCommand(ctx context.Context) *cobra.Command {
	var sceneName string
	var before string
	var after string
	var verbose bool

	cmd := cobra.Command{
		Use: "hotspots",
		Args: cobra.ExactArgs(1),
		Short: "Get the hotspots of an application in JSON formatted for d3.js",
		RunE: func(cmd *cobra.Command, args []string) error {
			ui := ctx.Ui
			ui.SetVerbose(verbose)
			ui.Log("Getting app hotspots...")

			connection, err := ctx.GetConnection()
			if err != nil {
				return err
			}
			defer connection.Close()

			c := NewCore(connection)
			a, beforeTime, afterTime, err := core.ExtractDateRangeAndAppFromArgs(connection, sceneName, args[0], before, after)
			if err != nil {
				return errors.Wrap(err, "Command failed")
			}

			hotspots, err := c.QueryHotspots(a, beforeTime, afterTime)

			ui.Ok()

			str, _ := json.MarshalIndent(hotspots, "", "  ")
			ui.Print(string(str))

			return nil
		},
	}

	cmd.Flags().StringVarP(&sceneName, "scene", "s", "", "Scene name")
	cmd.Flags().StringVarP(&before, "before", "a", "", "Fetch all the hotspots before this day")
	cmd.Flags().StringVarP(&after, "after", "b", "", "Fetch all the hotspots after this day")
	cmd.Flags().BoolVar(&verbose, "verbose", false, "display the log information")
	return &cmd
}

func NewRevisionTrends(ctx context.Context) *cobra.Command {
	var sceneName string
	var boundaryName string
	var before string
	var after string
	var csv bool
	var verbose bool

	cmd := cobra.Command{
		Use: "revision-trends",
		Aliases: []string{"revisions-trends", "rt"},
		Short: "Get the revision trends for a boundary",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ui := ctx.Ui
			ui.SetVerbose(verbose)

			connection, err := ctx.GetConnection()
			if err != nil {
				return err
			}
			defer connection.Close()

			c := NewCore(connection)
			a, beforeTime, afterTime, err := core.ExtractDateRangeAndAppFromArgs(connection, sceneName, args[0], before, after)
			if err != nil {
				return errors.Wrap(err, "Command failed")
			}

			b, err := c.boundary.QueryByAppIdAndName(a.Id, boundaryName)
			if err != nil {
				return errors.Wrap(err, "Boundary not found")
			}

			ui.Log("Getting revisions trends...")
			trends, err := c.RevisionTrends(a.Id, b, beforeTime, afterTime)
			if err != nil {
				return errors.Wrap(err, "Cannot get revisions trends")
			}

			ui.Ok()

			headers := []string{"date"}
			for _, t := range b.Transformations {
				headers = append(headers, t.Name)
			}
			table := ui.Table(headers, csv)
			for _, rt := range trends {
				cols := []string{rt.Date}
				for _, t := range b.Transformations {
					cols = append(cols, strconv.Itoa(rt.FindEntityRevision(t.Name).NumberOfRevisions))
				}
				table.Add(cols...)
			}
			table.Print()

			return nil
		},
	}

	cmd.Flags().StringVarP(&sceneName, "scene", "s", "", "Scene name")
	cmd.Flags().StringVarP(&before, "before", "a", "", "Fetch all the hotspots before this day")
	cmd.Flags().StringVarP(&after, "after", "b", "", "Fetch all the hotspots after this day")
	cmd.Flags().StringVarP(&boundaryName, "boundary", "", "", "Boundary to use")
	cmd.Flags().BoolVar(&csv, "csv", false, "get the results in csv format")
	cmd.Flags().BoolVar(&verbose, "verbose", false, "display the log information")

	return &cmd
}