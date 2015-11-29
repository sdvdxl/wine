package time

import "time"

const (
	DateFormat     = "2006-01-02"
	DateTimeFormat = "2006-01-02 15:04:05"
)

func ParseDate(val string) *time.Time {
	t, err := time.Parse(DateFormat, val)
	if err != nil {
		return nil
	}

	return &t
}

func ParseDateTime(val string) *time.Time {
	t, err := time.Parse(DateTimeFormat, val)
	if err != nil {
		return nil
	}

	return &t
}
