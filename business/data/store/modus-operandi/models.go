package modus_operandi

type WordCount struct {
	Word  string `json:"word" db:"word"`
	Count int    `json:"count" db:"count"`
}
