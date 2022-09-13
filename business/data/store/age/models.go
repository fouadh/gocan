package age

type EntityAge struct {
	Name string `db:"name" json:"name"`
	Age  int    `db:"age" json:"age"`
}
