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
