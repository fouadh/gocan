package revisions

import (
	context "com.fha.gocan/internal/platform"
	"github.com/spf13/cobra"
)

func NewCommand(ctx *context.Context) *cobra.Command {
	var sceneName string

	cmd := cobra.Command{
		Use: "revisions",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ui := ctx.Ui
			ui.Say("Getting app revisions...")
			ui.Say("file1")
			ui.Ok()
			return nil
		},
	}

	cmd.Flags().StringVarP(&sceneName, "scene", "s", "", "Scene name")
	return &cmd
}

