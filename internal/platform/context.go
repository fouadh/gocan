package context

import (
	"com.fha.gocan/internal/platform/db"
	"com.fha.gocan/internal/platform/terminal"
	"github.com/go-playground/validator"
)

type Context struct {
	Ui terminal.UI
	DataSource db.DataSource
	Validator *validator.Validate
}

func New(dsn string, ui terminal.UI) *Context {
	dataSource := db.SqlxDataSource{
		Dsn: dsn,
		Ui: ui,
	}

	return &Context{
		Ui:         ui,
		DataSource: &dataSource,
		Validator: validator.New(),
	}
}