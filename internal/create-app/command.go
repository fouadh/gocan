package create_app

import (
	context "com.fha.gocan/internal/platform"
	"github.com/spf13/cobra"
)

func NewCommand(ctx *context.Context) *cobra.Command {
	var sceneName string

	cmd := cobra.Command{
		Use: "create-app",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	cmd.Flags().StringVarP(&sceneName, "scene", "s", "", "Scene name")
	return &cmd
}
