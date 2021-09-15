package churn

import (
	"com.fha.gocan/business/core/app"
	"com.fha.gocan/foundation"
	"com.fha.gocan/foundation/date"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func NewCodeChurn(ctx foundation.Context) *cobra.Command {
	var sceneName string
	var before string
	var after string

	cmd := cobra.Command{
		Use: "code-churn",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ui := ctx.Ui
			connection, err := ctx.GetConnection()
			if err != nil {
				return err
			}

			ui.Say("Retrieving code churn...")

			c := NewCore(connection)

			beforeTime, err := date.ParseDay(before)
			if err != nil {
				return errors.Wrap(err, "Invalid before date")
			}

			afterTime, err := date.ParseDay(after)
			if err != nil {
				return errors.Wrap(err, "Invalid after date")
			}

			a, err := app.FindAppBySceneNameAndAppName(connection, sceneName, args[0])
			if err != nil {
				return errors.Wrap(err, "Command failed")
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
	cmd.Flags().StringVarP(&before, "before", "", date.Today(), "Fetch the code churn before this day")
	cmd.Flags().StringVarP(&after, "after", "", date.LongTimeAgo(), "Fetch all the code churn after this day")

	return &cmd
}