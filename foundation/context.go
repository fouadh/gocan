package foundation

import (
	"com.fha.gocan/foundation/db"
	"com.fha.gocan/foundation/terminal"
	"fmt"
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
		return nil, errors.Wrap(err, fmt.Sprintf("The connection to the dabase could not be established: %v", err.Error()))
	}

	return connection, nil
}