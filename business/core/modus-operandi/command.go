package modus_operandi

import (
	"com.fha.gocan/business/core"
	"com.fha.gocan/foundation"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"strconv"
)

func NewModusOperandi(ctx foundation.Context) *cobra.Command {
	var sceneName string
	var before string
	var after string
	var csv bool

	cmd := cobra.Command{
		Use: "modus-operandi",
		Args: cobra.ExactArgs(1),
		Short: "Get the most used terms in git messages of an application",
		RunE: func(cmd *cobra.Command, args []string) error {
			ui := ctx.Ui
			connection, err := ctx.GetConnection()
			if err != nil {
				return err
			}
			defer connection.Close()

			ui.Say("Retrieving modus operandi...")

			c := NewCore(connection)

			a, beforeTime, afterTime, err := core.ExtractDateRangeAndAppFromArgs(connection, sceneName, args[0], before, after)
			if err != nil {
				return errors.Wrap(err, "Invalid argument(s)")
			}

			words, err := c.Query(a.Id, beforeTime, afterTime)
			if err != nil {
				return err
			}

			table := ui.Table([]string{"word", "count"}, csv)
			for _, w := range words {
				table.Add(w.Word, strconv.Itoa(w.Count))
			}
			table.Print()

			return nil
		},
	}

	cmd.Flags().StringVarP(&sceneName, "scene", "s", "", "Scene name")
	cmd.Flags().StringVarP(&before, "before", "", "", "Fetch the modus-operandi before this day")
	cmd.Flags().StringVarP(&after, "after", "", "", "Fetch the modus-operandi after this day")
	cmd.Flags().BoolVar(&csv, "csv", false, "get the results in csv format")

	return &cmd
}