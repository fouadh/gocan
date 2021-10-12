package scene

import (
	"com.fha.gocan/business/data/store/scene"
	"com.fha.gocan/foundation"
	"com.fha.gocan/foundation/terminal"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func NewCreateScene(ctx foundation.Context) *cobra.Command {
	return &cobra.Command{
		Use:  "create-scene",
		Short: "Create a scene",
		Example: "gocan create-scene myscene",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ui := ctx.Ui
			ui.Say("Creating the scene...")

			connection, err := ctx.GetConnection()
			if err != nil {
				return err
			}
			defer connection.Close()

			core := NewCore(connection)
			s, err := core.Create(ctx, args[0])
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("Unable to create the scene: %s", err.Error()))
			}

			ui.Say(fmt.Sprintln("Scene", s.Id, "created."))
			ui.Ok()
			return nil
		},
	}
}

func NewScenes(ctx foundation.Context) *cobra.Command {
	var csv bool

	cmd := cobra.Command{
		Use:     "scenes",
		Short:   "List the scenes",
		Example: "gocan scenes",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ui := ctx.Ui

			connection, err := ctx.GetConnection()
			if err != nil {
				return err
			}
			defer connection.Close()

			ui.Say("Retrieving scenes...")
			core := NewCore(connection)
			scenes, err := core.QueryAll()
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("Unable to retrieve the scenes: %s", err.Error()))
			}

			if len(scenes) == 0 {
				ui.Say("No scene found")
				return nil
			}
			ui.Ok()

			printScenes(ui, scenes, csv)

			return nil
		},
	}

	cmd.Flags().BoolVar(&csv, "csv", false, "get the results in csv format")
	return &cmd
}

func NewDeleteScene(ctx foundation.Context) *cobra.Command {
	return &cobra.Command{
		Use: "delete-scene",
		Short: "Delete the specified scene",
		Args: cobra.ExactArgs(1),
		Example: "gocan delete-scene myscene",
		RunE: func(cmd *cobra.Command, args []string) error {
			ui := ctx.Ui

			connection, err := ctx.GetConnection()
			if err != nil {
				return err
			}
			defer connection.Close()

			ui.Say("Deleting scene...")
			c := NewCore(connection)

			if err := c.DeleteSceneByName(args[0]); err != nil {
				return errors.Wrap(err, "Unable to delete the scene")
			}

			ui.Ok()

			return nil
		},
	}
}

func printScenes(ui terminal.UI, scenes []scene.Scene, csv bool) {
	table := ui.Table([]string{
		"id",
		"name",
	}, csv)

	for _, s := range scenes {
		table.Add(s.Id, s.Name)
	}

	table.Print()
}
