package create_scene

import (
	context "com.fha.gocan/business/platform"
	"fmt"
	"github.com/pborman/uuid"
	"github.com/pkg/errors"
)

type CreateSceneRequest struct {
	Name string `validate:"required,max=255"`
}

func CreateScene(request CreateSceneRequest, ctx *context.Context) error {
	if err := ctx.Validator.Struct(request); err != nil {
		return errors.Wrap(err, "Invalid request to create a scene")
	}

	datasource := ctx.DataSource
	connection, err := datasource.GetConnection()
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("The connection to the dabase could not be established: %v", err.Error()))
	}

	id := uuid.NewUUID().String()
	_, err = connection.Exec("insert into scenes(id, name) values($1, $2)", id, request.Name)
	if err != nil {
		return errors.Wrap(err, "Scene could not be created")
	} else {
		return err
	}

	return nil
}
