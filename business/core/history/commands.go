package history

import (
	"com.fha.gocan/business/core"
	"com.fha.gocan/foundation"
	"com.fha.gocan/foundation/date"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func NewImportHistoryCommand(ctx *foundation.Context) *cobra.Command {
	var sceneName string
	var path string
	var before string
	var after string

	cmd := cobra.Command{
		Use:  "import-history",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ui := ctx.Ui
			connection, err := ctx.GetConnection()
			if err != nil {
				return err
			}

			ui.Say("Importing history...")

			a, beforeTime, afterTime, err := core.ExtractDateRangeAndAppFromArgs(connection, sceneName, args[0], before, after)
			if err != nil {
				return errors.Wrap(err, "Command failed")
			}

			h := NewCore(connection)
			if err = h.Import(a.Id, path, beforeTime, afterTime); err != nil {
				return errors.Wrap(err, "History cannot be imported")
			}

			ui.Ok()
			return nil
		},
	}

	cmd.Flags().StringVarP(&sceneName, "scene", "s", "", "Scene name")
	cmd.Flags().StringVarP(&path, "path", "d", ".", "App directory")
	cmd.Flags().StringVarP(&before, "before", "a", date.Today(), "Fetch all the hotspots before this day")
	cmd.Flags().StringVarP(&after, "after", "b", date.LongTimeAgo(), "Fetch all the hotspots after this day")
	return &cmd
}
