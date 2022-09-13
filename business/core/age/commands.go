package age

import (
	"com.fha.gocan/business/core"
	"com.fha.gocan/foundation"
	"com.fha.gocan/foundation/date"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func Commands(ctx foundation.Context) []*cobra.Command {
	return []*cobra.Command{
		get(ctx),
	}
}

func get(ctx foundation.Context) *cobra.Command {
	var sceneName string
	var before string
	var after string
	var initialDate string
	var csv bool
	var verbose bool

	cmd := cobra.Command{
		Use:   "code-age",
		Short: "Retrieve the age of the entities in months",
		Args:  cobra.ExactArgs(1),
		Example: `
gocan age myapp --scene myscene
gocan age myapp --scene myscene --after 2022-01-01 --before 2022-06-30
gocan age myapp --scene myscene --initial-date 2021-01-01
`,
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

			data, err := c.GetCodeAge(a.Id, initialDate, beforeTime, afterTime)
			if err != nil {
				return errors.Wrap(err, "Cannot calculate the code age")
			}

			ctx.Ui.Ok()

			table := ctx.Ui.Table([]string{"entity", "age"}, csv)

			for _, ea := range data {
				table.Add(ea.Name, fmt.Sprintf("%d", ea.Age))
			}

			table.Print()

			return nil
		},
	}

	cmd.Flags().StringVarP(&sceneName, "scene", "s", "", "Scene name")
	cmd.Flags().StringVarP(&before, "before", "", "", "Calculate the code age before this day")
	cmd.Flags().StringVarP(&after, "after", "", "", "Calculate the code age after this day")
	cmd.Flags().StringVarP(&initialDate, "initial-date", "", date.Today(), "From when to calculate the age")
	cmd.Flags().BoolVar(&csv, "csv", false, "get the results in csv format")

	cmd.MarkFlagRequired("scene")

	return &cmd
}
