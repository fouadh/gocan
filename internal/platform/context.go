package context

import (
	"com.fha.gocan/internal/platform/config"
	"com.fha.gocan/internal/platform/db"
	"com.fha.gocan/internal/platform/terminal"
	"github.com/go-playground/validator"
	"github.com/pkg/errors"
)

type Context struct {
	Ui terminal.UI
	DataSource db.DataSource
	Validator *validator.Validate
	Config *config.Config
}

func New(ui terminal.UI) (*Context, error) {
	config, err := config.ReadConfig()
	if err  != nil {
		return nil, errors.Wrap(err, "Could not read the configuration file. Please use gocan setup-db to eventually regenerate it")
	}

	dataSource := db.SqlxDataSource{
		Dsn: config.Dsn(),
		Ui: ui,
	}

	return &Context{
		Ui:         ui,
		DataSource: &dataSource,
		Validator: validator.New(),
		Config: config,
	}, nil
}