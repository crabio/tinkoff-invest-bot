package date

import (
	"time"
)

// GenerateDaySequence - generate sequence of days from start date to end
func GenerateDaySequence(startDate time.Time, endDate time.Time) (dateSequence []time.Time) {
	// Generate date sequence
	for date := startDate; date.After(endDate) == false; date = date.AddDate(0, 0, 1) {
		dateSequence = append(dateSequence, date)
	}

	return
}
