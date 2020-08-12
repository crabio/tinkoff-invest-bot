package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

// ExecQueryWithAttempts pollling query untill it's success or timeout is reached. 1 attempt = 1 second
func ExecQueryWithAttempts(db *sql.DB, queryStr string, maxAttempts uint) (err error) {
	// Until not success ot attemps out of range
	var success bool = false
	var attempt uint = 0
	for !success {
		log.Println("Execute query: ", queryStr)
		// Execute query
		_, err = db.Exec(queryStr)
		// Check err
		if err != nil {
			// Check attempts counter
			if attempt == maxAttempts {
				return err
			}
			// New attempt
			time.Sleep(1 * time.Second)
			attempt++
		} else {
			// Success request
			success = true
		}
	}
	return nil
}

// WaitDbInit pollling DB untill it's correctrly inited or timeout is reached. 1 attempt = 1 second
func WaitDbInit(config Configuration, maxAttempts uint) (err error) {
	log.Printf("Wait for DB init.")

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

	// Check instruments table
	err = ExecQueryWithAttempts(db, "SELECT 1 FROM instrument;", maxAttempts)
	// Check err
	if err != nil {
		log.Printf("Table instrument in DB is not inited.")
		return err
	}

	// Check candle intervals table
	err = ExecQueryWithAttempts(db, "SELECT 1 FROM candle_interval;", maxAttempts)
	// Check err
	if err != nil {
		log.Printf("Table candle_interval in DB is not inited.")
		return err
	}

	// Check candle table
	err = ExecQueryWithAttempts(db, "SELECT 1 FROM candle;", maxAttempts)
	// Check err
	if err != nil {
		log.Printf("Table candle in DB is not inited.")
		return err
	}

	return nil
}

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
	queryStr := `SELECT MAX(ts) FROM candle WHERE ts < $1;`

	// Execute query
	row := db.QueryRow(queryStr, startDate)

	// Convert timestamp from query to buffer
	err = row.Scan(&lattestTimestamp)

	return
}
