package db

// Configuration - struct which contains config info for Data Base
type Configuration struct {
	DbType                   string // Data Base type, for example "postgres"
	User                     string
	Password                 string
	Hosname                  string
	Port                     uint
	DbName                   string
	InstrumentsTableName     string
	CandleIntervalsTableName string
	CandlesTableName         string
}
