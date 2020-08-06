package config

import (
	"time"
)

// Configuration - struct which contains config info for app
type Configuration struct {
	SandboxToken          string
	ProductionToken       string
	GlobalRankCsvFilePath string
	StartLoadDate         time.Time
	MaxAttempts           uint
}
