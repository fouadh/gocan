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
	ComplexityId string    `db:"complexity_analysis_id" json:"complexityId,omitempty"`
	Lines        int       `db:"lines" json:"lines"`
	Indentations int       `db:"indentations" json:"indentations"`
	Mean         float64   `db:"mean" json:"mean"`
	Max          int       `db:"max" json:"max"`
	Stdev        float64   `db:"stdev" json:"stdev"`
	Date         time.Time `db:"date" json:"date"`
}

type ComplexityAnalysisSummary struct {
	Id   string `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}
