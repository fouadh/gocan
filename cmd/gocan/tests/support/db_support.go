package support

import (
	"com.fha.gocan/internal/platform"
	"com.fha.gocan/internal/platform/db"
)

func CreateDatabase(ctx *context.Context) *db.EmbeddedDatabase {
	database := db.EmbeddedDatabase{Config: ctx.Config}
	database.Start(ctx.Ui)
	db.Migrate(ctx.Config.Dsn(), ctx.Ui)
	return &database
}

