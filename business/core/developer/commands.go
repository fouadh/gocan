package developer

import (
	"com.fha.gocan/business/core"
	"com.fha.gocan/foundation"
	"com.fha.gocan/foundation/date"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"strconv"
)

func NewMainDevelopers(ctx foundation.Context) *cobra.Command {
	var sceneName string
	var before string
	var after string

	cmd := cobra.Command{
		Use:     "main-developers",
		Aliases: []string{"main_developers", "mainDevelopers", "md", "main-devs", "main_devs", "mainDevs"},
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ui := ctx.Ui
			connection, err := ctx.GetConnection()
			if err != nil {
				return err
			}

			ui.Say("Retrieving main developers...")

			c := NewCore(connection)

			a, beforeTime, afterTime, err := core.ExtractDateRangeAndAppFromArgs(connection, sceneName, args[0], before, after)
			if err != nil {
				return errors.Wrap(err, "Invalid argument(s)")
			}

			data, err := c.QueryMainDevelopers(a.Id, beforeTime, afterTime)
			if err != nil {
				return errors.Wrap(err, "Cannot retrieve main developers")
			}

			ui.Ok()

			table := ui.Table([]string{"entity", "main-dev", "added", "total-added", "ownership"})
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

func NewEntityEfforts(ctx foundation.Context) *cobra.Command {
	var sceneName string
	var before string
	var after string

	cmd := cobra.Command{
		Use:     "entity-efforts",
		Aliases: []string{"entity_efforts", "entityEfforts", "ee"},
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ui := ctx.Ui
			connection, err := ctx.GetConnection()
			if err != nil {
				return err
			}

			ui.Say("Retrieving entity efforts...")

			c := NewCore(connection)

			a, beforeTime, afterTime, err := core.ExtractDateRangeAndAppFromArgs(connection, sceneName, args[0], before, after)
			if err != nil {
				return errors.Wrap(err, "Invalid argument(s)")
			}

			data, err := c.QueryEntityEfforts(a.Id, beforeTime, afterTime)
			if err != nil {
				return errors.Wrap(err, "Cannot retrieve main developers")
			}

			ui.Ok()

			table := ui.Table([]string{"entity", "author", "author-revs", "total-revs"})
			for _, dev := range data {
				table.Add(dev.Entity, dev.Author, fmt.Sprint(dev.AuthorRevisions), fmt.Sprint(dev.TotalRevisions))
			}
			table.Print()

			return nil
		},
	}

	cmd.Flags().StringVarP(&sceneName, "scene", "s", "", "Scene name")
	cmd.Flags().StringVarP(&before, "before", "", date.Today(), "Fetch the entity efforts before this day")
	cmd.Flags().StringVarP(&after, "after", "", date.LongTimeAgo(), "Fetch all the entity efforts after this day")

	return &cmd
}

func NewKnowledgeMapCommand(ctx foundation.Context) *cobra.Command {
	var sceneName string
	var before string
	var after string

	cmd := cobra.Command{
		Use:     "knowledge-map",
		Aliases: []string{"knowledge_map", "knowledgeMap", "km"},
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ui := ctx.Ui
			connection, err := ctx.GetConnection()
			if err != nil {
				return err
			}

			ui.Say("Building knowledge map...")

			c := NewCore(connection)

			a, beforeTime, afterTime, err := core.ExtractDateRangeAndAppFromArgs(connection, sceneName, args[0], before, after)
			if err != nil {
				return errors.Wrap(err, "Invalid argument(s)")
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

func NewDevsCommand(ctx foundation.Context) *cobra.Command {
	var sceneName string
	var before string
	var after string

	cmd := cobra.Command{
		Use:     "devs",
		Aliases: []string{"developers"},
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ui := ctx.Ui
			ui.Say("Getting app revisions...")

			connection, err := ctx.GetConnection()
			if err != nil {
				return err
			}

			c := NewCore(connection)
			a, beforeTime, afterTime, err := core.ExtractDateRangeAndAppFromArgs(connection, sceneName, args[0], before, after)

			devs, err := c.QueryDevelopers(a.Id, beforeTime, afterTime)

			if err != nil {
				return err
			}
			ui.Ok()

			table := ui.Table([]string{"name", "commits"})
			for _, dev := range devs {
				table.Add(dev.Name, strconv.Itoa(dev.NumberOfCommits))
			}
			table.Print()

			return nil
		},
	}

	cmd.Flags().StringVarP(&sceneName, "scene", "s", "", "Scene name")
	cmd.Flags().StringVarP(&before, "before", "", date.Today(), "Fetch the developers before this day")
	cmd.Flags().StringVarP(&after, "after", "", date.LongTimeAgo(), "Fetch all the developers after this day")
	return &cmd
}
