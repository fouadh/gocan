package file_metrics

type FileContent struct {
	Language string `db:"language"`
	Files    string `db:"files"`
	Blank    string `db:"blanks"`
	Comment  string `db:"comments"`
	Code     string `db:"code"`
}
