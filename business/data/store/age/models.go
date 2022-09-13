package age

type EntityAge struct {
	Name string `db:"name" json:"name"`
	Age  int    `db:"age" json:"age"`
}

type EntityAgeHierarchy struct {
	Name     string                `json:"name"`
	Children []*EntityAgeHierarchy `json:"children,omitempty"`
	Weight   float64               `json:"weight,omitempty"`
	Size     int                   `json:"size,omitempty"`
}
