package revision

import (
	"com.fha.gocan/business/core"
	context "com.fha.gocan/foundation"
	"com.fha.gocan/foundation/date"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"strconv"
)

func NewRevisionsCommand(ctx context.Context) *cobra.Command {
	var sceneName string
	var before string
	var after string

	cmd := cobra.Command{
		Use:  "revisions",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ui := ctx.Ui
			ui.Say("Getting app revisions...")

			connection, err := ctx.GetConnection()
			if err != nil {
				return err
			}

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
			})
			for _, revision := range revisions {
				table.Add(revision.Entity, fmt.Sprint(revision.NumberOfRevisions))
			}
			table.Print()
			return nil
		},
	}

	cmd.Flags().StringVarP(&sceneName, "scene", "s", "", "Scene name")
	cmd.Flags().StringVarP(&before, "before", "a", date.Today(), "Fetch all the revisions before this day")
	cmd.Flags().StringVarP(&after, "after", "b", date.LongTimeAgo(), "Fetch all the revisions after this day")
	return &cmd
}

func NewHotspotsCommand(ctx context.Context) *cobra.Command {
	var sceneName string
	var before string
	var after string

	cmd := cobra.Command{
		Use: "hotspots",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ui := ctx.Ui
			ui.Say("Getting app hotspots...")

			connection, err := ctx.GetConnection()
			if err != nil {
				return err
			}

			c := NewCore(connection)
			a, beforeTime, afterTime, err := core.ExtractDateRangeAndAppFromArgs(connection, sceneName, args[0], before, after)
			if err != nil {
				return errors.Wrap(err, "Command failed")
			}

			hotspots, err := c.QueryHotspots(a, beforeTime, afterTime)

			ui.Ok()

			str, _ := json.MarshalIndent(hotspots, "", "  ")
			ui.Say(string(str))

			return nil
		},
	}

	cmd.Flags().StringVarP(&sceneName, "scene", "s", "", "Scene name")
	cmd.Flags().StringVarP(&before, "before", "a", date.Today(), "Fetch all the hotspots before this day")
	cmd.Flags().StringVarP(&after, "after", "b", date.LongTimeAgo(), "Fetch all the hotspots after this day")
	return &cmd
}

func NewRevisionTrends(ctx context.Context) *cobra.Command {
	var sceneName string
	var boundaryName string
	var before string
	var after string

	cmd := cobra.Command{
		Use: "revision-trends",
		Aliases: []string{"revisions-trends", "rt"},
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ui := ctx.Ui

			connection, err := ctx.GetConnection()
			if err != nil {
				return err
			}

			c := NewCore(connection)
			a, beforeTime, afterTime, err := core.ExtractDateRangeAndAppFromArgs(connection, sceneName, args[0], before, after)
			if err != nil {
				return errors.Wrap(err, "Command failed")
			}

			b, err := c.boundary.QueryByAppIdAndName(a.Id, boundaryName)
			if err != nil {
				return errors.Wrap(err, "Boundary not found")
			}

			ui.Say("Getting revisions trends...")
			trends, err := c.RevisionTrends(a.Id, b, beforeTime, afterTime)
			if err != nil {
				return errors.Wrap(err, "Cannot get revisions trends")
			}

			ui.Ok()

			headers := []string{"date"}
			for _, t := range b.Transformations {
				headers = append(headers, t.Name)
			}
			table := ui.Table(headers)
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
	cmd.Flags().StringVarP(&before, "before", "a", date.Today(), "Fetch all the hotspots before this day")
	cmd.Flags().StringVarP(&after, "after", "b", date.LongTimeAgo(), "Fetch all the hotspots after this day")
	cmd.Flags().StringVarP(&boundaryName, "boundary", "", "", "Boundary to use")
	return &cmd
}