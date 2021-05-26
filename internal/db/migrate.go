package db

import (
	"com.fha.gocan/internal/terminal"
	"embed"
	"fmt"
	"github.com/Boostport/migration"
	"github.com/Boostport/migration/driver/postgres"
	"github.com/jmoiron/sqlx"
	"os"
)

//go:embed migrations
var migrations embed.FS

func Migrate(dsn string, ui terminal.UI) *sqlx.DB {
	embedSource := &migration.EmbedMigrationSource{
		EmbedFS: migrations,
		Dir:     "migrations",
	}

	driver, err := postgres.New(dsn)
	if err != nil {
		ui.Failed(fmt.Sprintf("%s Cannot connect to the database: %+v\n", err))
		os.Exit(1)
	}
	applied, err := migration.Migrate(driver, embedSource, migration.Up, 0)

	ui.Say(fmt.Sprintf("Applied %d migrations", applied))

	return nil
}


