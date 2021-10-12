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
	var host string
	var port int
	var user string
	var password string
	var dbName string
	var externalDb bool
	var embeddedPath string

	cmd := &cobra.Command{
		Use:   "setup-db",
		Short: "Configure the database options. Caution: if you are changing the embedded db properties, you will lose all the existing data.",
		RunE: func(cmd *cobra.Command, args []string) error {
			ui := ctx.Ui
			ui.Log("Configuring the database...")
			var dataPath string
			if externalDb {
				dataPath = ""
			} else {
				dataPath = embeddedPath
			}

			c := db.Config{
				Host:             host,
				Port:             port,
				User:             user,
				Password:         password,
				Database:         dbName,
				Embedded:         !externalDb,
				EmbeddedDataPath: dataPath,
			}

			data, err := json.Marshal(c)
			if err != nil {
				return errors.Wrap(err, "Unable to marshal configuration object into json")
			}

			if err := createDirectory(embeddedPath); err != nil {
				return errors.Wrap(err, "Unable to create embedded db directory")
			}

			if err := createDirectory(defaultPath()); err != nil {
				return errors.Wrap(err, "Unable to create configuration directory")
			}

			if err := ioutil.WriteFile(defaultPath()+"/config.json", data, 0644); err != nil {
				return errors.Wrap(err, "Unable to save the configuration file")
			}

			ui.Log("Database configured")
			ui.Ok()
			return nil
		},
	}

	cmd.Flags().StringVarP(&host, "host", "", db.DefaultConfig.Host, "Database host")
	cmd.Flags().IntVarP(&port, "port", "", db.DefaultConfig.Port, "Port to which the database listens.")
	cmd.Flags().StringVarP(&user, "user", "u", db.DefaultConfig.User, "Database user")
	cmd.Flags().StringVarP(&password, "password", "p", db.DefaultConfig.Password, "Database password")
	cmd.Flags().StringVarP(&dbName, "database", "n", db.DefaultConfig.Database, "Database name")
	cmd.Flags().BoolVarP(&externalDb, "external-db", "e", false, "Set this flag if you prefer to use an external db rather than the embedded one")
	cmd.Flags().StringVarP(&embeddedPath, "directory", "d", defaultPath()+"/data", "Directory where the data will be stored. Only valid for the embedded database.")

	return cmd
}

func createDirectory(directory string) error {
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		if err := os.Mkdir(directory, os.ModeDir|0755); err != nil {
			return err
		}
	}
	return nil
}

func NewStartDbCommand(ctx *foundation.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start-db",
		Args: cobra.NoArgs,
		Short: "Start an embedded database",
		RunE: func(cmd *cobra.Command, args []string) error {
			if !ctx.Config.Embedded {
				ctx.Ui.Failed("The configuration specifies that an external database will be used: please use gocan setup-db if you want to switch to the embedded database.")
				return nil
			}
			ui := ctx.Ui
			database := db.EmbeddedDatabase{Config: ctx.Config}
			database.Start(ui)
			ui.Log("Applying migrations...")
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
		Args: cobra.NoArgs,
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

func NewMigrateDb(ctx *foundation.Context) *cobra.Command {
	cmd := cobra.Command{
		Use: "migrate-db",
		Short: "Run the migration scripts against a database. To be used with external dbs.",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			schema.Migrate(ctx.Config.Dsn(), ctx.Ui)
			return nil
		},
	}

	return &cmd
}

func defaultPath() string {
	usr, err := user.Current()
	if err != nil {
		fmt.Println("Cannot get user information: the config path will be created within this location")
		return ".gocan"
	}

	return usr.HomeDir + "/.gocan"
}
