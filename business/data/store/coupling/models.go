package coupling

type Coupling struct {
	Entity           string
	Coupled          string
	Degree           float64
	AverageRevisions float64
}

type Soc struct {
	Entity string
	Soc    int
}

type CouplingHierarchy struct {
	Name     string               `json:"name"`
	Children []*CouplingHierarchy `json:"children,omitempty"`
	Coupling []string             `json:"coupling,omitempty"`
	Relations []CouplingRelation  `json:"relations,omitempty"`
}

type CouplingRelation struct {
	Coupled          string  `json:"coupled"`
	Degree           float64 `json:"degree"`
	AverageRevisions float64 `json:"averageRevisions"`
}

