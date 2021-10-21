package stat

type Stat struct {
	AppId      string `db:"app_id"`
	CommitId   string `db:"commit_id"`
	Insertions int    `db:"insertions"`
	Deletions  int    `db:"deletions"`
	File       string `db:"file"`
}

type StatInfo struct {
	CommitId   string `db:"commit_id"`
	File       string `db:"file"`
}