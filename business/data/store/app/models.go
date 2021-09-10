package app

type App struct {
	Id      string `db:"id"`
	Name    string `db:"name"`
	SceneId string `db:"scene_id"`
}

type NewApp struct {
	Name    string `validate:"required,max=255"`
	SceneId string `validate:"required"`
}
