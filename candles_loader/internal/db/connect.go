package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// CreateDbConnection creates instance of connection to DB based on configuration
func CreateDbConnection(config Configuration) (dbConnection *sql.DB, err error) {
	// Create connection string
	connectionString := fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=disable",
		config.Type, config.User, config.Password, config.Hostname, config.Port, config.DbName)

	// Connect to DB
	dbConnection, err = sql.Open(config.Type, connectionString)
	// Check err
	if err != nil {
		return nil, err
	}

	return
}
