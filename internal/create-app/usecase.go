package create_app

import (
	"com.fha.gocan/internal/platform"
	"fmt"
	"github.com/pborman/uuid"
	"github.com/pkg/errors"
)

type CreateAppRequest struct {
	Name      string `validate:"required"`
	SceneName string
}

func CreateApp(request CreateAppRequest, ctx *context.Context) error {
	datasource := ctx.DataSource
	connection, err := datasource.GetConnection()
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("The connection to the dabase could not be established: %v", err.Error()))
	}
	id := uuid.NewUUID().String()
	_, err = connection.Exec("insert into apps(id, name, scene_id) values($1, $2, (select id from scenes where name=$3))", id, request.Name, request.SceneName)
	if err != nil {
		return errors.Wrap(err, "App could not be created")
	} else {
		return err
	}
	return nil
}
