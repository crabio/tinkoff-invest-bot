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
	CandleInterval        string
	// DB Configuration
	DbType     string // Data Base type, for example "postgres"
	DbUser     string
	DbPassword string
	DbHosname  string
	DbPort     uint
	DbName     string
}
