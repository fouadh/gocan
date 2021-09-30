package churn

import (
	"com.fha.gocan/business/core"
	"com.fha.gocan/foundation"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func NewCodeChurn(ctx foundation.Context) *cobra.Command {
	var sceneName string
	var before string
	var after string

	cmd := cobra.Command{
		Use:  "code-churn",
		Short: "Get the code churn of an application",
		Example: `
gocan code-churn myapp --scene myscene --after 2021-01-01 --before 2021-02-01
gocan code-churn myapp --scene myscene --before 2021-02-01
gocan code-churn myapp --scene myscene --after 2021-01-01
gocan code-churn myapp --scene myscene
`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ui := ctx.Ui
			connection, err := ctx.GetConnection()
			if err != nil {
				return err
			}

			ui.Say("Retrieving code churn...")

			c := NewCore(connection)

			a, beforeTime, afterTime, err := core.ExtractDateRangeAndAppFromArgs(connection, sceneName, args[0], before, after)
			if err != nil {
				return errors.Wrap(err, "Invalid argument(s)")
			}

			data, err := c.QueryCodeChurn(a.Id, beforeTime, afterTime)

			if err != nil {
				return errors.Wrap(err, "Cannot retrieve code churn")
			}

			ui.Ok()

			table := ui.Table([]string{"date", "added", "deleted"})

			for _, cc := range data {
				table.Add(cc.Date, fmt.Sprint(cc.Added), fmt.Sprint(cc.Deleted))
			}

			table.Print()

			return nil
		},
	}

	cmd.Flags().StringVarP(&sceneName, "scene", "s", "", "Scene name")
	cmd.Flags().StringVarP(&before, "before", "", "", "Fetch the code churn before this day")
	cmd.Flags().StringVarP(&after, "after", "", "", "Fetch all the code churn after this day")

	return &cmd
}

