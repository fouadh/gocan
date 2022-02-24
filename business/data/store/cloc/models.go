package cloc

type Clocs struct {
	Files []FileInfo
}

type FileInfo struct {
	AppId      string `db:"app_id"`
	Language   string `db:"language"`
	Extension  string `db:"extension"`
	Filename   string `db:"filename"`
	Location   string `db:"file"`
	Lines      int    `db:"lines"`
	Code       int    `db:"code"`
	Comment    int    `db:"comment"`
	Blank      int    `db:"blank"`
	Complexity int    `db:"complexity"`
	Binary     bool   `db:"is_binary"`
	CommitId   string `db:"commit_id"`
}
