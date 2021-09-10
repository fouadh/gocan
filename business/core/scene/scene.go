package scene

import (
	"com.fha.gocan/business/data/store/scene"
	context "com.fha.gocan/foundation"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type Core struct {
	scene scene.Store
}

func NewCore(connection *sqlx.DB) Core {
	return Core {
		scene: scene.NewStore(connection),
	}
}

func (c Core) Create(ctx context.Context, newScene scene.NewScene) (scene.Scene, error) {
	s, err := c.scene.Create(ctx, newScene)
	if err != nil {
		return scene.Scene{}, errors.Wrap(err, "create")
	}

	return s, nil
}
