package date_test

import (
	"github.com/iakrevetkho/tinkoff-invest-bot/candles_loader/internal/date"
	"testing"
	"time"
)

func TestGenerateDaySequence(t *testing.T) {
	// Run function
	result := date.GenerateDaySequence(time.Date(2020, 2, 13, 13, 54, 33, 69, time.UTC), time.Date(2020, 2, 17, 13, 40, 33, 69, time.UTC))
	// Create expected result
	expected := []time.Time{time.Date(2020, 2, 13, 13, 54, 33, 69, time.UTC),
		time.Date(2020, 2, 14, 13, 54, 33, 69, time.UTC),
		time.Date(2020, 2, 15, 13, 54, 33, 69, time.UTC)}
	// Check result len
	if len(result) != len(expected) {
		t.Error("Expected len ", len(expected), ", got ", len(result))
	}
	// Checck each element
	for i := range result {
		if result[i] != expected[i] {
			t.Error("Expected el #", i, " ", expected[i], ", got ", result[i])
		}
	}
}
