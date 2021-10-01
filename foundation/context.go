package foundation

import (
	"com.fha.gocan/foundation/db"
	"com.fha.gocan/foundation/terminal"
	"github.com/go-playground/validator"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type Context struct {
	Ui terminal.UI
	DataSource db.DataSource
	Validator *validator.Validate
	Config *db.Config
}

func New(ui terminal.UI, config *db.Config) *Context {
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

func (c Context) GetConnection() (*sqlx.DB, error)  {
	connection, err := c.DataSource.GetConnection()
	if err != nil {
		return nil, errors.Wrap(err, "Unable to connect to the database. Did you start it ?")
	}

	return connection, nil
}