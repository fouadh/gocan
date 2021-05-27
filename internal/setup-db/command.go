package setup_db

import (
	context "com.fha.gocan/internal/platform"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"os/user"
)

const dsn = "host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable"

type Config struct {
	Dsn string `json:"dsn"`
}

func BuildSetupDbCmd(ctx *context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use: "setup-db",
		RunE: func(cmd *cobra.Command, args []string) error {
			ui := ctx.Ui
			ui.Say("Configuring the database...")
			config := Config{
				Dsn: dsn,
			}
			data, err := json.Marshal(&config)
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("Failed to marshal configuration object into json: %v", err))
			}

			usr, err := user.Current()
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("Failed to get current user info: %v", err))
			}

			path := usr.HomeDir + "/.gocan"
			if _, err := os.Stat(path); os.IsNotExist(err) {
				err = os.Mkdir(path, os.ModeDir | 0755)
				if err != nil {
					return errors.Wrap(err, fmt.Sprintf("Failed to create gocan directory: %v", err))
				}
			}

			if err := ioutil.WriteFile(path + "/config.json", data, 0644); err != nil {
				return errors.Wrap(err, fmt.Sprintf("Failed to save configuration: %v", err))
			}

			ui.Ok()
			return nil
		},
	}

	return cmd
}