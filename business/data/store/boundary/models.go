package boundary

type Boundary struct {
	Id      string   `db:"id" json:"id"`
	Name    string   `db:"name" json:"name"`
	AppId   string   `db:"app_id" json:"appId"`
	Modules []Module `json:"modules"`
}

type NewBoundary struct {
	Name    string
	AppId   string
	Modules []NewModule
}

type Module struct {
	BoundaryId string `db:"boundary_id" json:"boundary_id"`
	Name       string `db:"name" json:"name"`
	Path       string `db:"path" json:"path"`
}

type NewModule struct {
	Name string
	Path string
}
