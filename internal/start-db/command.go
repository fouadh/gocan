package start_db

import (
	"com.fha.gocan/internal/platform/db"
	"com.fha.gocan/internal/platform/terminal"
	"github.com/spf13/cobra"
)

const dsn = "host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable"

func BuildStartDbCmd(ui terminal.UI) *cobra.Command {
	cmd := &cobra.Command{
		Use: "start-db",
		Short: "Start en embedded database",
		RunE: func(cmd *cobra.Command, args []string) error {
			database := db.EmbeddedDatabase{}
			database.Start(ui)
			db.Migrate(dsn, ui)
			return nil
		},
	}

	return cmd
}