package active_set

import (
	"com.fha.gocan/business/core"
	"com.fha.gocan/foundation"
	"com.fha.gocan/foundation/date"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"strconv"
)

func NewActiveSet(ctx foundation.Context) *cobra.Command {
	var sceneName string
	var before string
	var after string

	cmd := cobra.Command{
		Use:  "active-set",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			connection, err := ctx.GetConnection()
			if err != nil {
				return err
			}

			ctx.Ui.Say("Retrieving the apps...")
			c := NewCore(connection)

			a, beforeTime, afterTime, err := core.ExtractDateRangeAndAppFromArgs(connection, sceneName, args[0], before, after)
			if err != nil {
				return errors.Wrap(err, "Invalid argument(s)")
			}

			as, err := c.Query(a.Id, beforeTime, afterTime)
			if err != nil {
				return err
			}

			ctx.Ui.Ok()

			table := ctx.Ui.Table([]string{"Date", "Opened", "Closed"})
			for _, item := range as {
				table.Add(item.Date, strconv.Itoa(item.Opened), strconv.Itoa(item.Closed))
			}
			table.Print()

			return nil
		},
	}

	cmd.Flags().StringVarP(&sceneName, "scene", "s", "", "Scene name")
	cmd.Flags().StringVarP(&before, "before", "", date.Today(), "Fetch the active set before this day")
	cmd.Flags().StringVarP(&after, "after", "", date.LongTimeAgo(), "Fetch active set after this day")
	return &cmd
}
