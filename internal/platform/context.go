package context

import (
	"com.fha.gocan/internal/platform/db"
	"com.fha.gocan/internal/platform/terminal"
)

type Context struct {
	Ui terminal.UI
	DataSource db.DataSource
}

func New(dsn string, ui terminal.UI) *Context {
	dataSource := db.SqlxDataSource{
		Dsn: dsn,
		Ui: ui,
	}

	return &Context{
		Ui:         ui,
		DataSource: &dataSource,
	}
}