package start_db

import (
	context "com.fha.gocan/internal/platform"
	"com.fha.gocan/internal/platform/db"
	"github.com/spf13/cobra"
)

const dsn = "host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable"

func BuildStartDbCmd(ctx *context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use: "start-db",
		Short: "Start en embedded database",
		RunE: func(cmd *cobra.Command, args []string) error {
			ui := ctx.Ui
			database := db.EmbeddedDatabase{}
			database.Start(ui)
			return db.Migrate(dsn, ui)
		},
	}

	return cmd
}