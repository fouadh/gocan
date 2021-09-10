package stop_db

import (
	context "com.fha.gocan/business/platform"
	"com.fha.gocan/business/platform/db"
	"github.com/spf13/cobra"
)

func NewCommand(ctx *context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use: "stop-db",
		Short: "Stop the embedded database",
		RunE: func(cmd *cobra.Command, args []string) error {
			ui := ctx.Ui
			database := db.EmbeddedDatabase{Config: ctx.Config}
			database.Stop(ui)
			return nil
		},
	}

	return cmd
}