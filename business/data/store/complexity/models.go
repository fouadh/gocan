package complexity

import "time"

type Complexity struct {
	Id      string
	Name    string
	Entity  string
	AppId   string
	Entries []ComplexityEntry
}

type ComplexityEntry struct {
	Lines        int
	Indentations int
	Mean         float64
	Max          int
	Stdev        float64
	Date         time.Time
}
