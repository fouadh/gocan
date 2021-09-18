package active_set

import "time"

type ActiveSet struct {
	Date   time.Time `json:"date"`
	Opened int       `json:"opened,omitempty"`
	Closed int       `json:"closed,omitempty"`
}

type ActiveSetStats struct {
	Date  time.Time
	Count int
}
