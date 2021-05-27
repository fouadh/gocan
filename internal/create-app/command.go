package create_app

import (
	context "com.fha.gocan/internal/platform"
	"fmt"
	"github.com/pborman/uuid"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func NewCommand(ctx *context.Context) *cobra.Command {
	var sceneName string

	cmd := cobra.Command{
		Use: "create-app",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			datasource := ctx.DataSource
			connection, err := datasource.GetConnection()
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("The connection to the dabase could not be established: %v", err.Error()))
			}
			id := uuid.NewUUID().String()
			_, err = connection.Exec("insert into apps(id, name, scene_id) values($1, $2, (select id from scenes where name=$3))", id, args[0], sceneName)
			if err != nil {
				return errors.Wrap(err, "App could not be created")
			} else {
				return err
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&sceneName, "scene", "s", "", "Scene name")
	return &cmd
}
