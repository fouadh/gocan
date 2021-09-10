package start_db

import (
	"com.fha.gocan/business/data"
	context "com.fha.gocan/business/platform"
	"com.fha.gocan/business/platform/db"
	"github.com/spf13/cobra"
)

const dsn = "host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable"

func NewCommand(ctx *context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use: "start-db",
		Short: "Start en embedded database",
		RunE: func(cmd *cobra.Command, args []string) error {
			ui := ctx.Ui
			database := db.EmbeddedDatabase{Config: ctx.Config}
			database.Start(ui)
			ui.Say("Applying migrations...")
			data.Migrate(ctx.Config.Dsn(), ui)
			ui.Ok()
			return nil
		},
	}

	return cmd
}