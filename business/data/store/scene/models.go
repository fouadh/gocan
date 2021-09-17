package scene

type Scene struct {
	Id   string `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

type NewScene struct {
	Name string `validate:"required,max=255"`
}
