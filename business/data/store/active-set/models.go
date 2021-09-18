package active_set

import "time"

type ActiveSet struct {
	Date   string `json:"date"`
	Opened int       `json:"opened"`
	Closed int       `json:"closed"`
}

type ActiveSetStats struct {
	Date  time.Time
	Count int
}
