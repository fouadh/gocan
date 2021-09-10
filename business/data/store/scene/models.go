package scene

type Scene struct {
	Id   string `db:"id"`
	Name string `db:"name"`
}

type NewScene struct {
	Name string `validate:"required,max=255"`
}
