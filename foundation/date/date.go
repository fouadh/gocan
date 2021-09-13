package date

import "time"

func ParseDay(day string) (time.Time, error) {
	return time.Parse("2006-01-02", day)
}

func Today() string {
	return time.Now().AddDate(0, 0, 1).Format("2006-01-02")
}