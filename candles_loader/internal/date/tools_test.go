package date_test

import (
	"github.com/iakrevetkho/tinkoff-invest-bot/candles_loader/internal/date"
	"testing"
	"time"
)

func TestBeginOfDay(t *testing.T) {
	// Run function
	result := date.BeginOfDay(time.Date(2020, 2, 13, 13, 54, 33, 69, time.UTC))
	// Create expected result
	expected := time.Date(2020, 2, 13, 0, 0, 0, 0, time.UTC)
	// Check result
	if result != expected {
		t.Error("Expected ", expected, ", got ", result)
	}
}
