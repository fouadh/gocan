package coupling

import (
	"com.fha.gocan/business/core"
	"com.fha.gocan/foundation"
	"com.fha.gocan/foundation/date"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"math"
)

func NewCouplingCommand(ctx *foundation.Context) *cobra.Command {
	var sceneName string
	var minCoupling int
	var minRevsAvg int
	var before string
	var after string

	cmd := cobra.Command{
		Use:  "coupling",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ui := ctx.Ui
			connection, err := ctx.GetConnection()
			if err != nil {
				return err
			}

			ui.Say("Retrieving couplings...")

			c := NewCore(connection)

			a, beforeTime, afterTime, err := core.ExtractDateRangeAndAppFromArgs(connection, sceneName, args[0], before, after)
			if err != nil {
				return errors.Wrap(err, "Invalid argument(s)")
			}

			data, err := c.Query(a.Id, beforeTime, afterTime, float64(minCoupling)/100, minRevsAvg)

			if err != nil {
				return errors.Wrap(err, "Cannot retrieve couplings")
			}

			ui.Ok()

			if len(data) == 0 {
				ui.Say("No coupling found.")
				return nil
			}

			table := ui.Table([]string{"entity", "coupled", "degree", "average-revs"})
			for _, coupling := range data {
				table.Add(coupling.Entity, coupling.Coupled, fmt.Sprintf("%.2f", coupling.Degree), fmt.Sprint(int(math.Ceil(coupling.AverageRevisions))))
			}
			table.Print()

			return nil
		},
	}

	cmd.Flags().StringVarP(&sceneName, "scene", "s", "", "Scene name")
	cmd.Flags().IntVarP(&minCoupling, "min-degree", "d", 30, "minimal degree of coupling wanted (in percent)")
	cmd.Flags().IntVarP(&minRevsAvg, "min-revisions-average", "r", 5, "minimal number of average revisions wanted (in percent)")
	cmd.Flags().StringVarP(&before, "before", "", date.Today(), "Fetch all the couplings before this day")
	cmd.Flags().StringVarP(&after, "after", "", date.LongTimeAgo(), "Fetch all the couplings after this day")

	return &cmd
}

func NewSocCommand(ctx foundation.Context) *cobra.Command {
	var sceneName string
	var before string
	var after string

	cmd := cobra.Command{
		Use:  "soc",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ui := ctx.Ui
			connection, err := ctx.GetConnection()
			if err != nil {
				return err
			}

			ui.Say("Retrieving summary...")

			c := NewCore(connection)

			a, beforeTime, afterTime, err := core.ExtractDateRangeAndAppFromArgs(connection, sceneName, args[0], before, after)
			if err != nil {
				return errors.Wrap(err, "Invalid argument(s)")
			}

			data, err := c.QuerySoc(a.Id, beforeTime, afterTime)

			if err != nil {
				return errors.Wrap(err, "Cannot retrieve summary of coupling")
			}

			ui.Ok()

			table := ui.Table([]string { "entity", "soc" })
			for _, soc := range data {
				table.Add(soc.Entity, fmt.Sprint(soc.Soc))
			}
			table.Print()

			return nil
		},
	}

	cmd.Flags().StringVarP(&sceneName, "scene", "s", "", "Scene name")
	cmd.Flags().StringVarP(&before, "before", "", date.Today(), "Fetch the summary of coupling before this day")
	cmd.Flags().StringVarP(&after, "after", "", date.LongTimeAgo(), "Fetch all the summary of coupling after this day")

	return &cmd
}