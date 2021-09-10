package create_scene

import (
  context "com.fha.gocan/foundation"
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
      request := CreateSceneRequest{Name: args[0]}
      if err := CreateScene(request, ctx); err != nil {
        return errors.Wrap(err, "Error while creating the scene")
      }
      ui.Ok()
      return nil
    },
  }
}
