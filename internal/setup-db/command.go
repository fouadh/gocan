package setup_db

import (
	"com.fha.gocan/internal/platform/terminal"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"os/user"
)

const dsn = "host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable"

type Config struct {
	Dsn string `json:"dsn"`
}

func BuildSetupDbCmd(ui terminal.UI) *cobra.Command {
	cmd := &cobra.Command{
		Use: "setup-db",
		RunE: func(cmd *cobra.Command, args []string) error {
			ui.Say("Configuring the database...")
			config := Config{
				Dsn: dsn,
			}
			data, err := json.Marshal(&config)
			if err != nil {
				ui.Failed(fmt.Sprintf("Failed to marshal configuration object into json: %v", err))
				os.Exit(1)
			}

			usr, err := user.Current()
			if err != nil {
				ui.Failed(fmt.Sprintf("Failed to get current user info: %v", err))
				os.Exit(2)
			}

			path := usr.HomeDir + "/.gocan"
			if _, err := os.Stat(path); os.IsNotExist(err) {
				err = os.Mkdir(path, os.ModeDir | 0755)
				if err != nil {
					ui.Failed(fmt.Sprintf("Failed to create gocan directory: %v", err))
					os.Exit(3)
				}
			}

			if err := ioutil.WriteFile(path + "/config.json", data, 0644); err != nil {
				ui.Failed(fmt.Sprintf("Failed to save configuration: %v", err))
			}

			ui.Ok()
			return nil
		},
	}

	return cmd
}