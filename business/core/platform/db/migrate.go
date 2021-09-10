package db

import (
	terminal2 "com.fha.gocan/business/core/platform/terminal"
	"embed"
	"fmt"
	"github.com/Boostport/migration"
	"github.com/Boostport/migration/driver/postgres"
	"github.com/pkg/errors"
)

//go:embed migrations
var migrations embed.FS

func Migrate(dsn string, ui terminal2.UI) error {
	embedSource := &migration.EmbedMigrationSource{
		EmbedFS: migrations,
		Dir:     "migrations",
	}

	driver, err := postgres.New(dsn)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Cannot connect to the database: %+v\n", err))
	}

	applied, err := migration.Migrate(driver, embedSource, migration.Up, 0)

	if err != nil {
		return errors.Wrap(err, "Migrations could not be applied")
	}

	ui.Say(fmt.Sprintf("Applied %d migrations", applied))
	return nil
}


