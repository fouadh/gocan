package active_set

import (
	"com.fha.gocan/business/core"
	"com.fha.gocan/foundation"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"strconv"
)

func Commands(ctx foundation.Context) []*cobra.Command {
	return []*cobra.Command{
		list(ctx),
	}
}

func list(ctx foundation.Context) *cobra.Command {
	var sceneName string
	var before string
	var after string
	var csv bool
	var verbose bool

	cmd := cobra.Command{
		Use:  "active-set",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx.Ui.SetVerbose(verbose)
			connection, err := ctx.GetConnection()
			if err != nil {
				return err
			}
			defer connection.Close()
			ctx.Ui.Log("Retrieving the apps...")
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

			table := ctx.Ui.Table([]string{"Date", "Opened", "Closed"}, csv)
			for _, item := range as {
				table.Add(item.Date, strconv.Itoa(item.Opened), strconv.Itoa(item.Closed))
			}
			table.Print()

			return nil
		},
	}

	cmd.Flags().StringVarP(&sceneName, "scene", "s", "", "Scene name")
	cmd.Flags().StringVarP(&before, "before", "", "", "Fetch the active set before this day")
	cmd.Flags().StringVarP(&after, "after", "", "", "Fetch active set after this day")
	cmd.Flags().BoolVar(&csv, "csv", false, "get the results in csv format")
	cmd.Flags().BoolVar(&verbose, "verbose", false, "display the log information")

	cmd.MarkFlagRequired("scene")

	return &cmd
}
