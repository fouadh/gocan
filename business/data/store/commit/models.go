package commit

import "time"

type Commit struct {
	Id      string
	Author  string
	Date    time.Time
	Message string
	AppId   string
}

type CommitRange struct {
	MinDate time.Time `db:"min_date"`
	MaxDate time.Time `db:"max_date"`
}

func (cr CommitRange) MaxDay() time.Time {
	return time.Date(cr.MaxDate.Year(), cr.MaxDate.Month(), cr.MaxDate.Day(), 0, 0, 0, 0, cr.MaxDate.Location()).AddDate(0, 0, 1)
}

func (cr CommitRange) MinDay() time.Time {
	return time.Date(cr.MinDate.Year(), cr.MinDate.Month(), cr.MinDate.Day(), 0, 0, 0, 0, cr.MinDate.Location())
}
