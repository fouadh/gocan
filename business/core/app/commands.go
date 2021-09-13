package app

import (
	"com.fha.gocan/foundation"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func NewCreateAppCommand(ctx *foundation.Context) *cobra.Command {
	var sceneName string

	cmd := cobra.Command{
		Use:  "create-app",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx.Ui.Say("Creating the app...")

			connection, err := ctx.GetConnection()
			if err != nil {
				return err
			}

			core := NewCore(connection)
			a, err := core.Create(*ctx, args[0], sceneName)

			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("Unable to create the app: %s", err.Error()))
			}

			ctx.Ui.Say(fmt.Sprintln("App", a.Id, "created."))
			ctx.Ui.Ok()
			return nil
		},
	}
	cmd.Flags().StringVarP(&sceneName, "scene", "s", "", "Scene name")
	return &cmd
}

