package scene

import (
	"com.fha.gocan/foundation"
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

			datasource := ctx.DataSource
			connection, err := datasource.GetConnection()

			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("The connection to the dabase could not be established: %v", err.Error()))
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
