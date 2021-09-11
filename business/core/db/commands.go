package db

import (
	"com.fha.gocan/business/data/schema"
	"com.fha.gocan/foundation"
	"com.fha.gocan/foundation/db"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"os/user"
)

func NewSetupDbCommand(ctx *foundation.Context) *cobra.Command {
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

			c := db.Config{
				Host:             db.DefaultConfig.Host,
				Port:             db.DefaultConfig.Port,
				User:             db.DefaultConfig.User,
				Password:         db.DefaultConfig.Password,
				Database:         db.DefaultConfig.Database,
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

//const dsn = "host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable"

func NewStartDbCommand(ctx *foundation.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start-db",
		Short: "Start en embedded database",
		RunE: func(cmd *cobra.Command, args []string) error {
			ui := ctx.Ui
			database := db.EmbeddedDatabase{Config: ctx.Config}
			database.Start(ui)
			ui.Say("Applying migrations...")
			schema.Migrate(ctx.Config.Dsn(), ui)
			ui.Ok()
			return nil
		},
	}

	return cmd
}

func NewStopDbCommand(ctx *foundation.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop-db",
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