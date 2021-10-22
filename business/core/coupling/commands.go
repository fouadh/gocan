package coupling

import (
	"com.fha.gocan/business/core"
	"com.fha.gocan/business/data/store/coupling"
	"com.fha.gocan/foundation"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"math"
)

func Commands(ctx foundation.Context) []*cobra.Command {
	return []*cobra.Command{
		list(ctx),
		sum(ctx),
		hierarchy(ctx),
	}
}

func list(ctx foundation.Context) *cobra.Command {
	var sceneName string
	var minCoupling int
	var minRevsAvg int
	var temporalPeriod int
	var boundaryName string
	var before string
	var after string
	var csv bool
	var verbose bool

	cmd := cobra.Command{
		Use:  "coupling",
		Args: cobra.ExactArgs(1),
		Short: "Get the coupling relationships between entities for a given application",
		Example: `
gocan coupling myapp --scene myscene --min-degree 30 --min-revisions-average 10
gocan coupling myapp --scene myscene --min-degree 30
gocan coupling myapp --scene myscene
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ui := ctx.Ui
			ui.SetVerbose(verbose)
			connection, err := ctx.GetConnection()
			if err != nil {
				return err
			}
			defer connection.Close()

			ui.Log("Retrieving couplings...")

			c := NewCore(connection)

			a, beforeTime, afterTime, err := core.ExtractDateRangeAndAppFromArgs(connection, sceneName, args[0], before, after)
			if err != nil {
				return errors.Wrap(err, "Invalid argument(s)")
			}

			var data []coupling.Coupling

			if boundaryName == "" {
				data, err = c.Query(a.Id, float64(minCoupling)/100, minRevsAvg, temporalPeriod, beforeTime, afterTime)
			} else {
				data, err = c.QueryByBoundary(a.Id, boundaryName, float64(minCoupling)/100, minRevsAvg, temporalPeriod, beforeTime, afterTime)
			}

			if err != nil {
				return errors.Wrap(err, "Cannot retrieve couplings")
			}

			ui.Ok()

			if len(data) == 0 {
				ui.Log("No coupling found.")
				return nil
			}

			table := ui.Table([]string{"entity", "coupled", "degree", "average-revs"}, csv)
			for _, coupling := range data {
				table.Add(coupling.Entity, coupling.Coupled, fmt.Sprintf("%.2f", coupling.Degree), fmt.Sprint(int(math.Ceil(coupling.AverageRevisions))))
			}
			table.Print()

			return nil
		},
	}

	cmd.Flags().StringVarP(&sceneName, "scene", "s", "", "Scene name")
	cmd.Flags().StringVarP(&boundaryName, "boundary", "", "", "Optional boundary name to get the analysis for a specific boundary")
	cmd.Flags().IntVarP(&minCoupling, "min-degree", "d", 30, "minimal degree of coupling wanted (in percent)")
	cmd.Flags().IntVarP(&minRevsAvg, "min-revisions-average", "r", 5, "minimal number of average revisions wanted (in percent)")
	cmd.Flags().IntVarP(&temporalPeriod, "temporal-period", "", 1, "number of days to treat commits within as a single change")
	cmd.Flags().StringVarP(&before, "before", "", "", "Fetch coupling before this day")
	cmd.Flags().StringVarP(&after, "after", "", "", "Fetch coupling after this day")
	cmd.Flags().BoolVar(&csv, "csv", false, "get the results in csv format")
	cmd.Flags().BoolVar(&verbose, "verbose", false, "display the log information")
	return &cmd
}

func sum(ctx foundation.Context) *cobra.Command {
	var sceneName string
	var before string
	var after string
	var csv bool
	var verbose bool

	cmd := cobra.Command{
		Use:  "sum-of-coupling",
		Aliases: []string{"soc"},
		Args: cobra.ExactArgs(1),
		Short: "Get a sum of coupling for an application",
		Example: `
gocan sum-of-coupling myapp --scene myscene --after 2021-01-01 --before 2021-02-01
gocan sum-of-coupling myapp --scene myscene --after 2021-01-01
gocan sum-of-coupling myapp --scene myscene
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ui := ctx.Ui
			ui.SetVerbose(verbose)
			connection, err := ctx.GetConnection()
			if err != nil {
				return err
			}
			defer connection.Close()

			ui.Log("Retrieving sum...")

			c := NewCore(connection)

			a, beforeTime, afterTime, err := core.ExtractDateRangeAndAppFromArgs(connection, sceneName, args[0], before, after)
			if err != nil {
				return errors.Wrap(err, "Invalid argument(s)")
			}

			data, err := c.QuerySoc(a.Id, beforeTime, afterTime)

			if err != nil {
				return errors.Wrap(err, "Cannot retrieve sum of coupling")
			}

			ui.Ok()

			table := ui.Table([]string{"entity", "soc"}, csv)
			for _, soc := range data {
				table.Add(soc.Entity, fmt.Sprint(soc.Soc))
			}
			table.Print()

			return nil
		},
	}

	cmd.Flags().StringVarP(&sceneName, "scene", "s", "", "Scene name")
	cmd.Flags().StringVarP(&before, "before", "", "", "Fetch the sum of coupling before this day")
	cmd.Flags().StringVarP(&after, "after", "", "", "Fetch all the sum of coupling after this day")
	cmd.Flags().BoolVar(&csv, "csv", false, "get the results in csv format")
	cmd.Flags().BoolVar(&verbose, "verbose", false, "display the log information")
	return &cmd
}

func hierarchy(ctx foundation.Context) *cobra.Command {
	var sceneName string
	var minCoupling int
	var minRevsAvg int
	var before string
	var after string
	var verbose bool

	cmd := cobra.Command{
		Use: "coupling-hierarchy",
		Args: cobra.ExactArgs(1),
		Short: "Get the coupling information about an app in JSON formatted to be used with d3.js",
		RunE: func(cmd *cobra.Command, args []string) error {
			ui := ctx.Ui
			ui.SetVerbose(verbose)
			connection, err := ctx.GetConnection()
			if err != nil {
				return err
			}
			defer connection.Close()

			ui.Log("Retrieving sum of coupling...")

			c := NewCore(connection)

			a, beforeTime, afterTime, err := core.ExtractDateRangeAndAppFromArgs(connection, sceneName, args[0], before, after)
			if err != nil {
				return errors.Wrap(err, "Invalid argument(s)")
			}

			ch, err := c.BuildCouplingHierarchy(a, float64(minCoupling)/100, minRevsAvg, beforeTime, afterTime)

			ui.Ok()

			str, _ := json.MarshalIndent(ch, "", "  ")
			ui.Print(string(str))

			return nil
		},
	}

	cmd.Flags().StringVarP(&sceneName, "scene", "s", "", "Scene name")
	cmd.Flags().IntVarP(&minCoupling, "min-degree", "d", 30, "minimal degree of coupling wanted (in percent)")
	cmd.Flags().IntVarP(&minRevsAvg, "min-revisions-average", "r", 5, "minimal number of average revisions wanted (in percent)")
	cmd.Flags().StringVarP(&before, "before", "", "", "Fetch coupling before this day")
	cmd.Flags().StringVarP(&after, "after", "", "", "Fetch coupling after this day")
	cmd.Flags().BoolVar(&verbose, "verbose", false, "display the log information")
	return &cmd
}