package revision

import (
	context "com.fha.gocan/foundation"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"time"
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
			beforeTime := time.Now().AddDate(0, 0, 1)
			if after == "" {
				after = "1970-01-01"
			}
			afterTime, _ := time.Parse("2006-01-02", after)

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
	cmd.Flags().StringVarP(&before, "before", "", "", "")
	cmd.Flags().StringVarP(&after, "after", "", "", "")
	return &cmd
}