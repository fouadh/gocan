package boundary

type Boundary struct {
	Id              string           `db:"id" json:"id"`
	Name            string           `db:"name" json:"name"`
	AppId           string           `db:"app_id" json:"appId"`
	Transformations []Transformation `json:"transformations"`
}

type NewBoundary struct {
	Name            string
	AppId           string
	Transformations []NewTransformation
}

type Transformation struct {
	BoundaryId string `db:"boundary_id" json:"boundary_id"`
	Name       string `db:"name" json:"name"`
	Path       string `db:"path" json:"path"`
}

type NewTransformation struct {
	Name string
	Path string
}
