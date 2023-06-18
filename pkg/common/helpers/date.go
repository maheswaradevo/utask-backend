package helpers

import (
	"time"
)

var (
	DateLayoutTimeZone = "2006-01-02T15:04:05+08:00"
)

var (
	Location, _ = time.LoadLocation("Asia/Makassar")
)

func ParseDate(date *string, layout string) *time.Time {
	if date == nil {
		return nil
	}
	t, err := time.ParseInLocation(layout, *date, Location)
	if err != nil {
		return nil
	}
	return &t
}
