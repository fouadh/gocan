package revision

import (
	context "com.fha.gocan/foundation"
	"com.fha.gocan/foundation/date"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func NewRevisionsCommand(ctx *context.Context) *cobra.Command {
	var sceneName string
	var before string
	var after string

	cmd := cobra.Command{
		Use:  "revisions",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ui := ctx.Ui
			ui.Say("Getting app revisions...")

			datasource := ctx.DataSource
			connection, err := datasource.GetConnection()
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("The connection to the dabase could not be established: %v", err.Error()))
			}

			appName := args[0]

			core := NewCore(connection)
			beforeTime, err := date.ParseDay(before)
			if err != nil {
				return errors.Wrap(err, "Invalid before date")
			}

			afterTime, err := date.ParseDay(after)
			if err != nil {
				return errors.Wrap(err, "Invalid after date")
			}

			revisions, err := core.GetRevisions(*ctx, appName, sceneName, beforeTime, afterTime)

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
	cmd.Flags().StringVarP(&after, "after", "b", "1970-01-01", "Fetch all the revisions after this day")
	return &cmd
}

