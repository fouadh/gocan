package stop_db

import (
	"com.fha.gocan/internal/platform/db"
	"com.fha.gocan/internal/platform/terminal"
	"github.com/spf13/cobra"
)

func BuildStopDbCmd(ui terminal.UI) *cobra.Command {
	cmd := &cobra.Command{
		Use: "stop-db",
		Short: "Stop the embedded database",
		RunE: func(cmd *cobra.Command, args []string) error {
			database := db.EmbeddedDatabase{}
			database.Stop(ui)
			return nil
		},
	}

	return cmd
}