package import_history

import (
	context "com.fha.gocan/internal/platform"
	"github.com/spf13/cobra"
)

func NewCommand(ctx *context.Context) *cobra.Command {
	var sceneName string
	var path string

	cmd := cobra.Command{
		Use: "import-history",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ui := ctx.Ui
			ui.Say("Importing history...")
			ui.Ok()
			return nil
		},
	}

	cmd.Flags().StringVarP(&sceneName, "scene", "s", "", "Scene name")
	cmd.Flags().StringVarP(&path, "path", "d", ".", "App directory")
	return &cmd
}

