package complexity

import (
	"com.fha.gocan/business/core"
	"com.fha.gocan/foundation"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"strconv"
)

func NewComplexityAnalysis(ctx foundation.Context) *cobra.Command {
	var sceneName string
	var before string
	var after string
	var filename string

	cmd := cobra.Command{
		Use: "complexity-analysis",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ui := ctx.Ui
			connection, err := ctx.GetConnection()
			if err != nil {
				return err
			}

			ui.Say("Retrieving summary...")

			c := NewCore()

			a, beforeTime, afterTime, err := core.ExtractDateRangeAndAppFromArgs(connection, sceneName, args[0], before, after)
			if err != nil {
				return errors.Wrap(err, "Invalid argument(s)")
			}

			data, err := c.Analyze(a.Id, beforeTime, afterTime, filename)

			if err != nil {
				return errors.Wrap(err, "Error when analyzing complexity")
			}

			table := ui.Table([]string{"Indentations"})
			table.Add(strconv.Itoa(data.Indentations))
			table.Print()

			return nil
		},
	}

	cmd.Flags().StringVarP(&sceneName, "scene", "s", "", "Scene name")
	cmd.Flags().StringVarP(&before, "before", "", "", "Analyze the complexity before this day")
	cmd.Flags().StringVarP(&after, "after", "", "", "Analyze the complexity after this day")
	cmd.Flags().StringVarP(&filename, "filename", "f", "", "The file to analyze")

	return &cmd
}
