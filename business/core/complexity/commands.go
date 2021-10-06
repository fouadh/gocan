package complexity

import (
	"com.fha.gocan/business/core"
	"com.fha.gocan/business/core/app"
	"com.fha.gocan/foundation"
	"com.fha.gocan/foundation/date"
	"github.com/dustin/go-humanize"
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
	var spaces int

	cmd := cobra.Command{
		Use:     "create-complexity-analysis",
		Aliases: []string{"cca"},
		Short:   "Create a complexity analysis",
		Long: `
A complexity analysis is built by computing the number of lines and indentations in
the specified entity.
`,
		Example: `
gocan create-complexity-analysis myanalysis --app myapp --scene myscene --directory /code/project/ --filename src/main/File.java --after 2021-01-01 --before 2021-02-01 --spaces 2
gocan create-complexity-analysis myanalysis --app myapp --scene myscene --directory /code/project/ --filename src/main/File.java --after 2021-01-01 --spaces 2
gocan create-complexity-analysis myanalysis --app myapp --scene myscene --directory /code/project/ --filename src/main/File.java --before 2021-02-01 --spaces 2
gocan create-complexity-analysis myanalysis --app myapp --scene myscene --directory /code/project/ --filename src/main/File.java --spaces 2
gocan create-complexity-analysis myanalysis --app myapp --scene myscene --directory /code/project/ --filename src/main/File.java
`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ui := ctx.Ui
			connection, err := ctx.GetConnection()
			if err != nil {
				return err
			}
			defer connection.Close()

			c := NewCore(connection)

			a, beforeTime, afterTime, err := core.ExtractDateRangeAndAppFromArgs(connection, sceneName, appName, before, after)
			if err != nil {
				return errors.Wrap(err, "Invalid argument(s)")
			}

			ui.Say("Analyzing the file revisions between " + date.FormatDay(afterTime) + " and " + date.FormatDay(beforeTime))

			data, err := c.CreateComplexityAnalysis(args[0], a.Id, beforeTime, afterTime, filename, directory, spaces)

			if err != nil {
				return errors.Wrap(err, "Error when analyzing complexity")
			}

			ui.Ok()

			table := ui.Table([]string{"Date", "Lines", "Indentations", "Mean", "Stdev", "Max"})
			for _, cy := range data.Entries {
				table.Add(cy.Date.String(), strconv.Itoa(cy.Lines),
					strconv.Itoa(cy.Indentations),
					humanize.FtoaWithDigits(cy.Mean, 2),
					humanize.FtoaWithDigits(cy.Stdev, 2),
					strconv.Itoa(cy.Max))
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
	cmd.Flags().IntVarP(&spaces, "spaces", "", 4, "The number of spaces defining an indentation")

	return &cmd
}

func NewDeleteComplexityAnalysis(ctx foundation.Context) *cobra.Command {
	var sceneName string
	var appName string

	cmd := cobra.Command{
		Use:     "delete-complexity-analysis",
		Aliases: []string{"dca"},
		Short: "Delete a complexity analysis",
		Example: "gocan delete-complexity-analysis myanalysis --app myapp --scene myscene",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ui := ctx.Ui
			connection, err := ctx.GetConnection()
			if err != nil {
				return err
			}
			defer connection.Close()

			c := NewCore(connection)

			a, err := app.FindAppByAppNameAndSceneName(connection, appName, sceneName)
			if err != nil {
				return errors.Wrap(err, "Unable to find the app")
			}

			ui.Say("Deleting the analysis...")

			if err := c.DeleteAnalysisByName(a.Id, args[0]); err != nil {
				return errors.Wrap(err, "Unable to delete the analysis")
			}

			ui.Ok()

			return nil
		},
	}

	cmd.Flags().StringVarP(&sceneName, "scene", "s", "", "Scene name")
	cmd.Flags().StringVarP(&appName, "app", "a", "", "Application name")

	return &cmd
}

func NewComplexityAnalyses(ctx foundation.Context) *cobra.Command {
	var sceneName string
	var appName string

	cmd := cobra.Command{
		Use:     "complexity-analyses",
		Aliases: []string{"ca"},
		Args: cobra.NoArgs,
		Short: "List the complexity analyses associated with an app",
		Example: "gocan complexity-analyses --app myapp --scene myscene",
		RunE: func(cmd *cobra.Command, args []string) error {
			ui := ctx.Ui
			connection, err := ctx.GetConnection()
			if err != nil {
				return err
			}
			defer connection.Close()

			c := NewCore(connection)

			a, err := app.FindAppByAppNameAndSceneName(connection, appName, sceneName)
			if err != nil {
				return errors.Wrap(err, "Unable to find the app")
			}

			ui.Say("Fetching the analyses...")

			data, err := c.QueryAnalyses(a.Id)
			if err != nil {
				return errors.Wrap(err, "Unable to fetch analyses")
			}

			if len(data) > 0 {
				table := ui.Table([]string{"id", "name"})
				for _, a := range data {
					table.Add(a.Id, a.Name)
				}
				table.Print()
			} else {
				ui.Say("No analysis found.")
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&sceneName, "scene", "s", "", "Scene name")
	cmd.Flags().StringVarP(&appName, "app", "a", "", "Application name")

	return &cmd
}
