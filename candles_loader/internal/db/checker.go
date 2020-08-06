package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

// GetLattestCandleTimestamp gets lattest candle timestamp from DB, starts from startDate
func GetLattestCandleTimestamp(config Configuration, startDate time.Time) (lattestTimestamp time.Time, err error) {
	// Create onnection string
	connectionString := fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=disable",
		config.Type, config.User, config.Password, config.Hosname, config.Port, config.DbName)

	// Connect to DB
	db, err := sql.Open(config.Type, connectionString)
	// Check err
	if err != nil {
		return
	}
	// At the end close connetion
	defer db.Close()

	// Get unknown Instruments
	log.Println("Get unknown instruments list")
	queryStr := `
	SELECT MAX(ts)
	FROM candle as candle
	WHERE ts > $1;`

	// Execute query
	row := db.QueryRow(queryStr, startDate)

	// Convert timestamp from query to buffer
	err = row.Scan(&lattestTimestamp)

	return
}
