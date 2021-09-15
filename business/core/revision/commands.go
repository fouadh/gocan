package revision

import (
	"com.fha.gocan/business/core/app"
	context "com.fha.gocan/foundation"
	"com.fha.gocan/foundation/date"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
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

			core := NewCore(connection)
			beforeTime, err := date.ParseDay(before)
			if err != nil {
				return errors.Wrap(err, "Invalid before date")
			}

			afterTime, err := date.ParseDay(after)
			if err != nil {
				return errors.Wrap(err, "Invalid after date")
			}

			a, err := app.FindAppBySceneNameAndAppName(connection, sceneName, args[0])
			if err != nil {
				return errors.Wrap(err, "Command failed")
			}

			revisions, err := core.GetRevisions(a.Id, beforeTime, afterTime)

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

			core := NewCore(connection)
			beforeTime, err := date.ParseDay(before)
			if err != nil {
				return errors.Wrap(err, "Invalid before date")
			}

			afterTime, err := date.ParseDay(after)
			if err != nil {
				return errors.Wrap(err, "Invalid after date")
			}

			a, err := app.FindAppBySceneNameAndAppName(connection, sceneName, args[0])
			if err != nil {
				return errors.Wrap(err, "Command failed")
			}

			hotspots, err := core.GetHotspots(a, beforeTime, afterTime)

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