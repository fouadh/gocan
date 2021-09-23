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
