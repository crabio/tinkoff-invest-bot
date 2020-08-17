package db

// Configuration - struct which contains config info for Data Base
type Configuration struct {
	Type     string // Data Base type, for example "postgres"
	User     string
	Password string
	Hosname  string
	Port     uint
	DbName   string
}
