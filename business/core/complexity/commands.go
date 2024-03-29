package complexity

import (
	"com.fha.gocan/business/core"
	"com.fha.gocan/business/core/app"
	boundary2 "com.fha.gocan/business/core/boundary"
	"com.fha.gocan/business/data/store/boundary"
	"com.fha.gocan/business/data/store/complexity"
	"com.fha.gocan/foundation"
	"com.fha.gocan/foundation/date"
	"com.fha.gocan/foundation/terminal"
	"github.com/dustin/go-humanize"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"strconv"
)

func Commands(ctx foundation.Context) []*cobra.Command {
	return []*cobra.Command{
		create(ctx),
		delete(ctx),
		list(ctx),
		get(ctx),
	}
}

func create(ctx foundation.Context) *cobra.Command {
	var sceneName string
	var appName string
	var before string
	var after string
	var filename string
	var directory string
	var spaces int
	var csv bool
	var verbose bool
	var boundaryName string
	var moduleName string

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
			ui.SetVerbose(verbose)
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

			ui.Log("Analyzing the file revisions between " + date.FormatDay(afterTime) + " and " + date.FormatDay(beforeTime))

			module := boundary.Module{}

			bc := boundary2.NewCore(connection)
			if boundaryName != "" {
				b, err := bc.QueryByAppIdAndName(a.Id, boundaryName)
				if err != nil {
					return errors.Wrap(err, "Boundary not found")
				}
				for _, t := range b.Modules {
					if t.Name == moduleName {
						module = t
					}
				}
			}

			data, err := c.CreateComplexityAnalysis(args[0], a.Id, beforeTime, afterTime, filename, directory, spaces, module)

			if err != nil {
				return errors.Wrap(err, "Error when analyzing complexity")
			}

			ui.Ok()

			printComplexityEntries(ui, csv, data.Entries)

			return nil
		},
	}

	cmd.Flags().StringVarP(&sceneName, "scene", "s", "", "Scene name")
	cmd.Flags().StringVarP(&boundaryName, "boundary", "", "", "Boundary name (to be used with module name)")
	cmd.Flags().StringVarP(&moduleName, "module", "t", "", "Module name")
	cmd.Flags().StringVarP(&appName, "app", "a", "", "Application name")
	cmd.Flags().StringVarP(&before, "before", "", "", "Analyze the complexity before this day")
	cmd.Flags().StringVarP(&after, "after", "", "", "Analyze the complexity after this day")
	cmd.Flags().StringVarP(&filename, "filename", "f", "", "The file to analyze relative to the directory argument")
	cmd.Flags().StringVarP(&directory, "directory", "d", "", "The directory of the git repo")
	cmd.Flags().IntVarP(&spaces, "spaces", "", 4, "The number of spaces defining an indentation")
	cmd.Flags().BoolVar(&csv, "csv", false, "get the results in csv format")
	cmd.Flags().BoolVar(&verbose, "verbose", false, "display the log information")

	cmd.MarkFlagRequired("scene")
	cmd.MarkFlagRequired("app")
	cmd.MarkFlagRequired("directory")

	return &cmd
}

func printComplexityEntries(ui terminal.UI, csv bool, entries []complexity.ComplexityEntry) {
	table := ui.Table([]string{"Date", "Lines", "Indentations", "Mean", "Stdev", "Max"}, csv)
	for _, cy := range entries {
		table.Add(cy.Date.String(), strconv.Itoa(cy.Lines),
			strconv.Itoa(cy.Indentations),
			humanize.FtoaWithDigits(cy.Mean, 2),
			humanize.FtoaWithDigits(cy.Stdev, 2),
			strconv.Itoa(cy.Max))
	}
	table.Print()
}

func delete(ctx foundation.Context) *cobra.Command {
	var sceneName string
	var appName string
	var verbose bool

	cmd := cobra.Command{
		Use:     "delete-complexity-analysis",
		Aliases: []string{"dca"},
		Short:   "Delete a complexity analysis",
		Example: "gocan delete-complexity-analysis myanalysis --app myapp --scene myscene",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ui := ctx.Ui
			ui.SetVerbose(verbose)
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

			ui.Log("Deleting the analysis...")

			if err := c.DeleteAnalysisByName(a.Id, args[0]); err != nil {
				return errors.Wrap(err, "Unable to delete the analysis")
			}

			ui.Print("Complexity has been deleted.")
			ui.Ok()

			return nil
		},
	}

	cmd.Flags().StringVarP(&sceneName, "scene", "s", "", "Scene name")
	cmd.Flags().StringVarP(&appName, "app", "a", "", "Application name")
	cmd.Flags().BoolVar(&verbose, "verbose", false, "display the log information")

	cmd.MarkFlagRequired("scene")
	cmd.MarkFlagRequired("app")

	return &cmd
}

func get(ctx foundation.Context) *cobra.Command {
	var sceneName string
	var appName string
	var csv bool
	var verbose bool

	cmd := cobra.Command{
		Use:     "complexity-analysis",
		Args:    cobra.ExactArgs(1),
		Short:   "Retrieve a complexity analysis by its name",
		Example: "gocan complexity-analysis myfile --app myapp --scene myscene",
		RunE: func(cmd *cobra.Command, args []string) error {
			ui := ctx.Ui
			ui.SetVerbose(verbose)
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

			analysisName := args[0]
			ui.Log("Retrieving analysis " + analysisName + "...")

			data, err := c.QueryAnalyses(a.Id)
			if err != nil {
				return errors.Wrap(err, "Unable to fetch analyses")
			}

			var found []complexity.ComplexityEntry = nil

			if len(data) > 0 {
				for _, an := range data {
					if an.Name == analysisName {
						found, err = c.QueryAnalysisEntriesById(an.Id)
						if err != nil {
							return errors.Wrap(err, "Error while looking for analysis "+an.Id)
						}
					}

				}
			} else {
				ui.Log("No analysis found.")
			}

			if found != nil {
				printComplexityEntries(ui, csv, found)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&sceneName, "scene", "s", "", "Scene name")
	cmd.Flags().StringVarP(&appName, "app", "a", "", "Application name")
	cmd.Flags().BoolVar(&csv, "csv", false, "get the results in csv format")
	cmd.Flags().BoolVar(&verbose, "verbose", false, "display the log information")

	cmd.MarkFlagRequired("scene")
	cmd.MarkFlagRequired("app")

	return &cmd
}

func list(ctx foundation.Context) *cobra.Command {
	var sceneName string
	var appName string
	var csv bool
	var verbose bool

	cmd := cobra.Command{
		Use:     "complexity-analyses",
		Aliases: []string{"ca"},
		Args:    cobra.NoArgs,
		Short:   "List the complexity analyses associated with an app",
		Example: "gocan complexity-analyses --app myapp --scene myscene",
		RunE: func(cmd *cobra.Command, args []string) error {
			ui := ctx.Ui
			ui.SetVerbose(verbose)
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

			ui.Log("Fetching the analyses...")

			data, err := c.QueryAnalyses(a.Id)
			if err != nil {
				return errors.Wrap(err, "Unable to fetch analyses")
			}

			if len(data) > 0 {
				table := ui.Table([]string{"id", "name"}, csv)
				for _, a := range data {
					table.Add(a.Id, a.Name)
				}
				table.Print()
			} else {
				ui.Log("No analysis found.")
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&sceneName, "scene", "s", "", "Scene name")
	cmd.Flags().StringVarP(&appName, "app", "a", "", "Application name")
	cmd.Flags().BoolVar(&csv, "csv", false, "get the results in csv format")
	cmd.Flags().BoolVar(&verbose, "verbose", false, "display the log information")

	cmd.MarkFlagRequired("scene")
	cmd.MarkFlagRequired("app")

	return &cmd
}
