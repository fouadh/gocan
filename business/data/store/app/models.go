package app

type App struct {
	Id      string `db:"id" json:"id"`
	Name    string `db:"name" json:"name"`
	SceneId string `db:"scene_id" json:"sceneId"`
}

type NewApp struct {
	Name    string `validate:"required,max=255"`
	SceneId string `validate:"required"`
}

type Summary struct {
	Id                      string `json:"id"`
	Name                    string `json:"name"`
	NumberOfCommits         int    `json:"numberOfCommits"`
	NumberOfEntities        int    `json:"numberOfEntities"`
	NumberOfEntitiesChanged int    `json:"numberOfEntitiesChanged"`
	NumberOfAuthors         int    `json:"numberOfAuthors"`
}
