package stop_db

import (
	context "com.fha.gocan/internal/platform"
	"com.fha.gocan/internal/platform/db"
	"github.com/spf13/cobra"
)

func BuildStopDbCmd(ctx *context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use: "stop-db",
		Short: "Stop the embedded database",
		RunE: func(cmd *cobra.Command, args []string) error {
			ui := ctx.Ui
			database := db.EmbeddedDatabase{}
			database.Stop(ui)
			return nil
		},
	}

	return cmd
}