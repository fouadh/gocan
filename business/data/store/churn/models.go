package churn

type CodeChurn struct {
	Date    string `json:"date"`
	Added   int    `json:"added"`
	Deleted int    `json:"deleted"`
}
