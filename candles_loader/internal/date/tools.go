package date

import (
	"time"
)

// BeginOfDay return start of day from timestamp
func BeginOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}
