package developer

import (
	"com.fha.gocan/business/core"
	"com.fha.gocan/foundation"
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
	var csv bool
	var verbose bool

	cmd := cobra.Command{
		Use:     "main-developers",
		Aliases: []string{"main_developers", "mainDevelopers", "md", "main-devs", "main_devs", "mainDevs"},
		Short: "Get the main developers of an application entities",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ui := ctx.Ui
			ui.SetVerbose(verbose)
			connection, err := ctx.GetConnection()
			if err != nil {
				return err
			}
			defer connection.Close()

			ui.Log("Retrieving main developers...")

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

			table := ui.Table([]string{"entity", "main-dev", "added", "total-added", "ownership"}, csv)
			for _, dev := range data {
				table.Add(dev.Entity, dev.Author, fmt.Sprint(dev.Added), fmt.Sprint(dev.TotalAdded), fmt.Sprintf("%.2f", dev.Ownership))
			}
			table.Print()

			return nil
		},
	}

	cmd.Flags().StringVarP(&sceneName, "scene", "s", "", "Scene name")
	cmd.Flags().StringVarP(&before, "before", "", "", "Fetch the main developers before this day")
	cmd.Flags().StringVarP(&after, "after", "", "", "Fetch all the main developers after this day")
	cmd.Flags().BoolVar(&csv, "csv", false, "get the results in csv format")
	cmd.Flags().BoolVar(&verbose, "verbose", false, "display the log information")

	return &cmd
}

func NewEntityEfforts(ctx foundation.Context) *cobra.Command {
	var sceneName string
	var before string
	var after string
	var csv bool
	var verbose bool

	cmd := cobra.Command{
		Use:     "entity-efforts",
		Aliases: []string{"entity_efforts", "entityEfforts", "ee"},
		Short: "Get the efforts associated with entities of an application",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ui := ctx.Ui
			ui.SetVerbose(verbose)
			connection, err := ctx.GetConnection()
			if err != nil {
				return err
			}
			defer connection.Close()

			ui.Log("Retrieving entity efforts...")

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

			table := ui.Table([]string{"entity", "author", "author-revs", "total-revs"}, csv)
			for _, dev := range data {
				table.Add(dev.Entity, dev.Author, fmt.Sprint(dev.AuthorRevisions), fmt.Sprint(dev.TotalRevisions))
			}
			table.Print()

			return nil
		},
	}

	cmd.Flags().StringVarP(&sceneName, "scene", "s", "", "Scene name")
	cmd.Flags().StringVarP(&before, "before", "", "", "Fetch the entity efforts before this day")
	cmd.Flags().StringVarP(&after, "after", "", "", "Fetch all the entity efforts after this day")
	cmd.Flags().BoolVar(&csv, "csv", false, "get the results in csv format")
	cmd.Flags().BoolVar(&verbose, "verbose", false, "display the log information")

	return &cmd
}

func NewKnowledgeMapCommand(ctx foundation.Context) *cobra.Command {
	var sceneName string
	var before string
	var after string
	var verbose bool

	cmd := cobra.Command{
		Use:     "knowledge-map",
		Aliases: []string{"knowledge_map", "knowledgeMap", "km"},
		Short: "Get the knowledge map",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ui := ctx.Ui
			ui.SetVerbose(verbose)
			connection, err := ctx.GetConnection()
			if err != nil {
				return err
			}
			defer connection.Close()

			ui.Log("Building knowledge map...")

			c := NewCore(connection)

			a, beforeTime, afterTime, err := core.ExtractDateRangeAndAppFromArgs(connection, sceneName, args[0], before, after)
			if err != nil {
				return errors.Wrap(err, "Invalid argument(s)")
			}

			km, err := c.BuildKnowledgeMap(a, beforeTime, afterTime)

			ui.Ok()

			str, _ := json.MarshalIndent(km, "", "  ")
			ui.Print(string(str))

			return nil
		},
	}

	cmd.Flags().StringVarP(&sceneName, "scene", "s", "", "Scene name")
	cmd.Flags().StringVarP(&before, "before", "", "", "Fetch the main developers before this day")
	cmd.Flags().StringVarP(&after, "after", "", "", "Fetch all the main developers after this day")
	cmd.Flags().BoolVar(&verbose, "verbose", false, "display the log information")

	return &cmd
}

func NewDevsCommand(ctx foundation.Context) *cobra.Command {
	var sceneName string
	var before string
	var after string
	var csv bool
	var verbose bool

	cmd := cobra.Command{
		Use:     "devs",
		Aliases: []string{"developers"},
		Short: "Get the developers of an application",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ui := ctx.Ui
			ui.SetVerbose(verbose)
			ui.Log("Getting app revisions...")

			connection, err := ctx.GetConnection()
			if err != nil {
				return err
			}
			defer connection.Close()

			c := NewCore(connection)
			a, beforeTime, afterTime, err := core.ExtractDateRangeAndAppFromArgs(connection, sceneName, args[0], before, after)

			devs, err := c.QueryDevelopers(a.Id, beforeTime, afterTime)

			if err != nil {
				return err
			}
			ui.Ok()

			table := ui.Table([]string{"name", "commits"}, csv)
			for _, dev := range devs {
				table.Add(dev.Name, strconv.Itoa(dev.NumberOfCommits))
			}
			table.Print()

			return nil
		},
	}

	cmd.Flags().StringVarP(&sceneName, "scene", "s", "", "Scene name")
	cmd.Flags().StringVarP(&before, "before", "", "", "Fetch the developers before this day")
	cmd.Flags().StringVarP(&after, "after", "", "", "Fetch all the developers after this day")
	cmd.Flags().BoolVar(&csv, "csv", false, "get the results in csv format")
	cmd.Flags().BoolVar(&verbose, "verbose", false, "display the log information")

	return &cmd
}
