package date

import "time"

func ParseDay(day string) (time.Time, error) {
	return time.Parse("2006-01-02", day)
}

func FormatDay(t time.Time) string {
	return t.Format("2006-01-02")
}


func Today() string {
	return time.Now().AddDate(0, 0, 1).Format("2006-01-02")
}

func LongTimeAgo() string {
	return time.Now().AddDate(-1, 0, 0).Format("2006-01-02")
}