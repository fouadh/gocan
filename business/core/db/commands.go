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
	"strconv"
)

func Commands(ctx foundation.Context) []*cobra.Command {
	return []*cobra.Command {
		setupDbCommand(ctx),
		startDbCommand(ctx),
		stopDbCommand(ctx),
		migrateDb(ctx),
	}
}

func setupDbCommand(ctx foundation.Context) *cobra.Command {
	var host string
	var port int
	var user string
	var password string
	var dbName string
	var externalDb bool
	var embeddedPath string
	var verbose bool

	cmd := &cobra.Command{
		Use:   "setup-db",
		Short: "Configure the database options. Caution: if you are changing the embedded db properties, you will lose all the existing data.",
		RunE: func(cmd *cobra.Command, args []string) error {
			ui := ctx.Ui
			ui.SetVerbose(verbose)
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

			if !externalDb {
				ui.Log("Creating embedded database folder at " + embeddedPath)
				if err := createDirectory(embeddedPath); err != nil {
					return errors.Wrap(err, "Unable to create embedded db directory")
				}
			}

			ui.Log("Create configuration folder if not existing")
			if err := createDirectory(defaultPath()); err != nil {
				return errors.Wrap(err, "Unable to create configuration directory")
			}

			configLocation := defaultPath() + "/config.json"
			ui.Log("Saving configuration file at location " + configLocation)
			if err := ioutil.WriteFile(configLocation, data, 0644); err != nil {
				return errors.Wrap(err, "Unable to save the configuration file")
			}

			ui.Print("Database has been configured with the following properties:")
			t := ui.Table([]string{"host", "port", "user", "database", "embedded", "embedded db path"}, false)
			var embedded string
			if c.Embedded {
				embedded = "true"
			} else {
				embedded = "false"
			}
			t.Add(c.Host, strconv.Itoa(c.Port), c.User, c.Database, embedded, c.EmbeddedDataPath)
			t.Print()
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
	cmd.Flags().BoolVar(&verbose, "verbose", false, "display the log information")

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

func startDbCommand(ctx foundation.Context) *cobra.Command {
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

func stopDbCommand(ctx foundation.Context) *cobra.Command {
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

func migrateDb(ctx foundation.Context) *cobra.Command {
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
