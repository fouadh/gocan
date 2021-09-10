package create_scene

import (
  "com.fha.gocan/business/core/scene"
  scene2 "com.fha.gocan/business/data/store/scene"
  context "com.fha.gocan/foundation"
  "fmt"
  "github.com/pkg/errors"
  "github.com/spf13/cobra"
)

func NewCommand(ctx *context.Context) *cobra.Command {
  return &cobra.Command{
    Use: "create-scene",
    Args:  cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
      ui := ctx.Ui
      ui.Say("Creating the scene...")
      newScene := scene2.NewScene{Name: args[0]}

      datasource := ctx.DataSource
      connection, err := datasource.GetConnection()

      if err != nil {
        return errors.Wrap(err, fmt.Sprintf("The connection to the dabase could not be established: %v", err.Error()))
      }

      core := scene.NewCore(connection)
      s, err := core.Create(*ctx, newScene)
      if err != nil {
        return errors.Wrap(err, fmt.Sprintf("Unable to create the scene: %s", err.Error()))
      }

      ui.Say(fmt.Sprintln("Scene", s.Id, "created."))
      ui.Ok()
      return nil
    },
  }
}
