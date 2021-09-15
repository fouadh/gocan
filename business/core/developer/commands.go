package developer

import (
	"com.fha.gocan/business/core/app"
	"com.fha.gocan/foundation"
	"com.fha.gocan/foundation/date"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func NewMainDevelopers(ctx foundation.Context) *cobra.Command {
	var sceneName string
	var before string
	var after string

	cmd := cobra.Command{
		Use: "main-developers",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ui := ctx.Ui
			connection, err := ctx.GetConnection()
			if err != nil {
				return err
			}

			ui.Say("Retrieving main developers...")

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

			data, err := c.QueryMainDevelopers(a.Id, beforeTime, afterTime)
			if err != nil {
				return errors.Wrap(err, "Cannot retrieve main developers")
			}

			ui.Ok()

			table := ui.Table([]string {"entity", "main-dev", "added", "total-added", "ownership"})
			for _, dev := range data {
				table.Add(dev.Entity, dev.Author, fmt.Sprint(dev.Added), fmt.Sprint(dev.TotalAdded), fmt.Sprintf("%.2f", dev.Ownership))
			}
			table.Print()

			return nil
		},
	}

	cmd.Flags().StringVarP(&sceneName, "scene", "s", "", "Scene name")
	cmd.Flags().StringVarP(&before, "before", "", date.Today(), "Fetch the main developers before this day")
	cmd.Flags().StringVarP(&after, "after", "", date.LongTimeAgo(), "Fetch all the main developers after this day")

	return &cmd
}

func NewKnowledgeMapCommand(ctx foundation.Context) *cobra.Command {
	var sceneName string
	var before string
	var after string

	cmd := cobra.Command{
		Use: "knowledge-map",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ui := ctx.Ui
			connection, err := ctx.GetConnection()
			if err != nil {
				return err
			}

			ui.Say("Building knowledge map...")

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

			km, err := c.BuildKnowledgeMap(a, beforeTime, afterTime)

			ui.Ok()

			str, _ := json.MarshalIndent(km, "", "  ")
			ui.Say(string(str))

			return nil
		},
	}

	cmd.Flags().StringVarP(&sceneName, "scene", "s", "", "Scene name")
	cmd.Flags().StringVarP(&before, "before", "", date.Today(), "Fetch the main developers before this day")
	cmd.Flags().StringVarP(&after, "after", "", date.LongTimeAgo(), "Fetch all the main developers after this day")

	return &cmd
}