package date

import "time"

func ParseDay(day string) (time.Time, error) {
	return time.Parse("2006-01-02", day)
}

func ParseDateTime(datetime string) (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05", datetime)
}

func FormatDay(t time.Time) string {
	return t.Format("2006-01-02")
}

func Today() string {
	return time.Now().Format("2006-01-02")
}

func OneYearAgo() string {
	return time.Now().AddDate(-1, 0, 0).Format("2006-01-02")
}
