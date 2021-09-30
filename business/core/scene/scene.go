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
	return Core{
		scene: scene.NewStore(connection),
	}
}

func (c Core) Create(ctx context.Context, sceneName string) (scene.Scene, error) {
	newScene := scene.NewScene{Name: sceneName}
	s, err := c.scene.Create(ctx, newScene)
	if err != nil {
		return scene.Scene{}, errors.Wrap(err, "create")
	}

	return s, nil
}

func (c Core) QueryAll() ([]scene.Scene, error) {
	return c.scene.QueryAll()
}

func (c Core) QueryById(id string) (scene.Scene, error) {
	return c.scene.QueryById(id)
}

func (c Core) DeleteSceneByName(name string) error {
	return c.scene.DeleteByName(name)
}
