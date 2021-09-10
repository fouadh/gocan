package create_app

import (
	context "com.fha.gocan/foundation"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func NewCommand(ctx *context.Context) *cobra.Command {
	var sceneName string

	cmd := cobra.Command{
		Use: "create-app",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx.Ui.Say("Creating the app...")
			request := CreateAppRequest{
				Name:      args[0],
				SceneName: sceneName,
			}

			if err := CreateApp(request, ctx); err != nil {
				return errors.Wrap(err, "Unable to create the app")
			}

			ctx.Ui.Ok()
			return nil
		},
	}
	cmd.Flags().StringVarP(&sceneName, "scene", "s", "", "Scene name")
	return &cmd
}

