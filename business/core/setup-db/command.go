package setup_db

import (
	context "com.fha.gocan/business/core/platform"
	"com.fha.gocan/business/core/platform/config"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"os/user"
)

func NewCommand(ctx *context.Context) *cobra.Command {
	var dataPath string

	cmd := &cobra.Command{
		Use: "setup-db",
		RunE: func(cmd *cobra.Command, args []string) error {
			ui := ctx.Ui
			ui.Say("Configuring the database...")
			usr, err := user.Current()
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("Failed to get current user info: %v", err))
			}

			path := usr.HomeDir + "/.gocan"
			if dataPath == "" {
				dataPath = path
			} else {
				ui.Say("Data path will be set to: " + dataPath)
			}

			c := config.Config{
				Host:             config.DefaultConfig.Host,
				Port:             config.DefaultConfig.Port,
				User:             config.DefaultConfig.User,
				Password:         config.DefaultConfig.Password,
				Database:         config.DefaultConfig.Database,
				Embedded:         true,
				EmbeddedDataPath: dataPath + "/data",
			}

			data, err := json.Marshal(c)
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("Failed to marshal configuration object into json: %v", err))
			}

			if _, err := os.Stat(path); os.IsNotExist(err) {
				err = os.Mkdir(path, os.ModeDir|0755)
				if err != nil {
					return errors.Wrap(err, fmt.Sprintf("Failed to create gocan directory: %v", err))
				}
			}

			if err := ioutil.WriteFile(path+"/config.json", data, 0644); err != nil {
				return errors.Wrap(err, fmt.Sprintf("Failed to save configuration: %v", err))
			}

			ui.Say("Database configured")
			ui.Ok()
			return nil
		},
	}

	cmd.Flags().StringVarP(&dataPath, "path", "p", "", "Path where the postgresql data will be stored")

	return cmd
}
