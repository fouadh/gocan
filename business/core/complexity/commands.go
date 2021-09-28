package complexity

import (
	"com.fha.gocan/business/core"
	"com.fha.gocan/foundation"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"strconv"
)

func NewCreateComplexityAnalysis(ctx foundation.Context) *cobra.Command {
	var sceneName string
	var appName string
	var before string
	var after string
	var filename string
	var directory string

	cmd := cobra.Command{
		Use:  "complexity-analysis",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ui := ctx.Ui
			connection, err := ctx.GetConnection()
			if err != nil {
				return err
			}

			c := NewCore()

			a, beforeTime, afterTime, err := core.ExtractDateRangeAndAppFromArgs(connection, sceneName, appName, before, after)
			if err != nil {
				return errors.Wrap(err, "Invalid argument(s)")
			}

			ui.Say("Analyzing the file revisions...")

			data, err := c.CreateComplexityAnalysis(args[0], a.Id, beforeTime, afterTime, filename, directory)

			if err != nil {
				return errors.Wrap(err, "Error when analyzing complexity")
			}

			table := ui.Table([]string{"Date", "Lines", "Indentations"})
			for _, cy := range data {
				table.Add(cy.Date.String(), strconv.Itoa(cy.Lines), strconv.Itoa(cy.Indentations))
			}
			table.Print()

			return nil
		},
	}

	cmd.Flags().StringVarP(&sceneName, "scene", "s", "", "Scene name")
	cmd.Flags().StringVarP(&appName, "app", "a", "", "Application name")
	cmd.Flags().StringVarP(&before, "before", "", "", "Analyze the complexity before this day")
	cmd.Flags().StringVarP(&after, "after", "", "", "Analyze the complexity after this day")
	cmd.Flags().StringVarP(&filename, "filename", "f", "", "The file to analyze relative to the directory argument")
	cmd.Flags().StringVarP(&directory, "directory", "d", "", "The directory of the git repo")

	return &cmd
}
