package history

import (
	"com.fha.gocan/business/core"
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
	var verbose bool

	cmd := cobra.Command{
		Use:  "import-history",
		Args: cobra.ExactArgs(1),
		Short: "Import the commits of an application",
		RunE: func(cmd *cobra.Command, args []string) error {
			ui := ctx.Ui
			ui.SetVerbose(verbose)
			connection, err := ctx.GetConnection()
			if err != nil {
				return err
			}
			defer connection.Close()

			a, beforeTime, afterTime, err := core.ExtractDateRangeAndAppFromArgs(connection, sceneName, args[0], before, after)
			if err != nil {
				return errors.Wrap(err, "Command failed")
			}

			h := NewCore(connection)

			if err := h.CheckIfCanImport(path); err != nil {
				return errors.Wrap(err, "Unable to import history")
			}


			ui.Log("Importing history...")
			if err = h.Import(a.Id, path, beforeTime, afterTime, ctx); err != nil {
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
	cmd.Flags().StringVarP(&after, "after", "", date.LongTimeAgo(), "Fetch all the hotspots after this day")
	cmd.Flags().BoolVar(&verbose, "verbose", false, "display the log information")
	return &cmd
}
