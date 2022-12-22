package history

import (
	"com.fha.gocan/business/core/app"
	"com.fha.gocan/foundation"
	"com.fha.gocan/foundation/date"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func Commands(ctx foundation.Context) []*cobra.Command {
	return []*cobra.Command{
		create(ctx),
	}
}

func create(ctx foundation.Context) *cobra.Command {
	var sceneName string
	var path string
	var before string
	var after string
	var beforeCommit string
	var afterCommit string
	var intervalBetweenAnalyses int
	var verbose bool

	cmd := cobra.Command{
		Use:   "import-history",
		Args:  cobra.ExactArgs(1),
		Short: "Import the commits of an application",
		RunE: func(cmd *cobra.Command, args []string) error {
			ui := ctx.Ui
			ui.SetVerbose(verbose)
			connection, err := ctx.GetConnection()
			if err != nil {
				return err
			}
			defer connection.Close()

			a, err := app.FindAppByAppNameAndSceneName(connection, args[0], sceneName)

			h := NewCore(connection)

			if err := h.CheckIfCanImport(path); err != nil {
				return errors.Wrap(err, "Unable to import history")
			}

			ui.Log("Importing history...")
			if err = h.Import(a.Id, path, before, after, beforeCommit, afterCommit, ctx, intervalBetweenAnalyses); err != nil {
				return errors.Wrap(err, "History cannot be imported")
			}

			ui.Print("History has been imported")
			ui.Ok()
			return nil
		},
	}

	cmd.Flags().StringVarP(&sceneName, "scene", "s", "", "Scene name")
	cmd.Flags().StringVarP(&path, "directory", "d", ".", "App directory")
	cmd.Flags().StringVarP(&before, "before", "", date.Today(), "Fetch all the history before this day")
	cmd.Flags().StringVarP(&after, "after", "", "", "Fetch all the history after this day")
	cmd.Flags().StringVarP(&beforeCommit, "beforeCommit", "", "", "Fetch all the history before this commit")
	cmd.Flags().StringVarP(&afterCommit, "afterCommit", "", "", "Fetch all the history after this commit")
	cmd.Flags().IntVarP(&intervalBetweenAnalyses, "interval-between-analyses", "", 0, "Number of commits between two complexity analyses. By default, only one is done for the last commit.")
	cmd.Flags().BoolVar(&verbose, "verbose", false, "display the log information")

	cmd.MarkFlagRequired("scene")

	return &cmd
}
