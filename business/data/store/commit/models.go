package commit

import "time"

type Commit struct {
	Id      string
	Author  string
	Date    time.Time
	Message string
	AppId   string
}

