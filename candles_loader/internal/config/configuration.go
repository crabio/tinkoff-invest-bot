package config

import (
	"time"
)

// ConfigurationFile - struct which contains config info for app from file
type ConfigurationFile struct {
	SandboxToken    string
	ProductionToken string
	StartLoadDate   time.Time
}

// ConfigurationEnv - struct which contains config info for app fron environment variables
type ConfigurationEnv struct {
	GlobalRankCsvFilePath string
	MaxAttempts           uint
	CandleInterval        string
	// DB Configuration
	DbType     string // Data Base type, for example "postgres"
	DbUser     string
	DbPassword string
	DbHostname  string
	DbPort     uint
	DbName     string
}
