package active_set

import "time"

type ActiveSet struct {
	Date   time.Time `json:"date"`
	Opened int       `json:"opened"`
	Closed int       `json:"closed"`
}

type ActiveSetStats struct {
	Date  time.Time
	Count int
}
