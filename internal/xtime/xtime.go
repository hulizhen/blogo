package xtime

import "time"

const (
	layoutShort = "2006-01-02"
	layoutLong  = "2006-01-02 15:04:05"
)

func ShortFormat(t time.Time) string {
	return t.Format(layoutShort)
}

func LongFormat(t time.Time) string {
	return t.Format(layoutLong)
}
