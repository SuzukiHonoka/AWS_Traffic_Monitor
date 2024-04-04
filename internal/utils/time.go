package utils

import "time"

func BeginningOfDay(now time.Time) time.Time {
	y, m, _ := now.Date()
	return time.Date(y, m, 1, 0, 0, 0, 0, time.UTC)
}

func EndOfDay(now time.Time) time.Time {
	y, m, d := now.Date()
	return time.Date(y, m, d, 23, 59, 59, int(time.Second-time.Nanosecond), time.UTC)
}
