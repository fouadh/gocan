package context

import (
	"com.fha.gocan/internal/platform/config"
	"com.fha.gocan/internal/platform/db"
	"com.fha.gocan/internal/platform/terminal"
	"github.com/go-playground/validator"
)

type Context struct {
	Ui terminal.UI
	DataSource db.DataSource
	Validator *validator.Validate
	Config *config.Config
}

func New(ui terminal.UI, config *config.Config) *Context {
	dataSource := db.SqlxDataSource{
		Dsn: config.Dsn(),
		Ui: ui,
	}

	return &Context{
		Ui:         ui,
		DataSource: &dataSource,
		Validator: validator.New(),
		Config: config,
	}
}