package context

import (
	"com.fha.gocan/business/platform/config"
	"com.fha.gocan/business/platform/db"
	"com.fha.gocan/foundation/terminal"
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