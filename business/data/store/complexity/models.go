package complexity

import "time"

type Complexity struct {
	Id      string `db:"id"`
	Name    string `db:"name"`
	Entity  string `db:"entity"`
	AppId   string `db:"app_id"`
	Entries []ComplexityEntry
}

type ComplexityEntry struct {
	ComplexityId string    `db:"complexity_analysis_id"`
	Lines        int       `db:"lines"`
	Indentations int       `db:"indentations"`
	Mean         float64   `db:"mean"`
	Max          int       `db:"max"`
	Stdev        float64   `db:"stdev"`
	Date         time.Time `db:"date"`
}
