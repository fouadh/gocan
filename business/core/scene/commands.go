package scene

import (
	"com.fha.gocan/business/data/store/scene"
	"com.fha.gocan/foundation"
	"com.fha.gocan/foundation/terminal"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func NewCreateSceneCommand(ctx *foundation.Context) *cobra.Command {
	return &cobra.Command{
		Use:  "create-scene",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ui := ctx.Ui
			ui.Say("Creating the scene...")

			connection, err := ctx.GetConnection()
			if err != nil {
				return err
			}

			core := NewCore(connection)
			s, err := core.Create(*ctx, args[0])
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("Unable to create the scene: %s", err.Error()))
			}

			ui.Say(fmt.Sprintln("Scene", s.Id, "created."))
			ui.Ok()
			return nil
		},
	}
}

func NewScenesCommand(ctx *foundation.Context) *cobra.Command {
	return &cobra.Command{
		Use: "scenes",
		RunE: func(cmd *cobra.Command, args []string) error {
			ui := ctx.Ui

			connection, err := ctx.GetConnection()
			if err != nil {
				return err
			}

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

			printScenes(ui, scenes)

			return nil
		},
	}
}

func printScenes(ui terminal.UI, scenes []scene.Scene) {
	table := ui.Table([]string{
		"name",
	})

	for _, s := range scenes {
		table.Add(s.Name)
	}

	table.Print()
}
