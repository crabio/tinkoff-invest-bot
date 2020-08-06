package db

import (
	"database/sql"
	"fmt"
	"github.com/iakrevetkho/tinkoff-invest-bot/candles_loader/internal/date"
	"log"
	"time"
)

// DeleteCandlesFromDay deletecandles from DB starts from specific day
func DeleteCandlesFromDay(config Configuration, timestamp time.Time) (err error) {
	// Create connection string
	connectionString := fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=disable",
		config.Type, config.User, config.Password, config.Hosname, config.Port, config.DbName)

	// Connect to DB
	db, err := sql.Open(config.Type, connectionString)
	// Check err
	if err != nil {
		return err
	}
	// At the end close connetion
	defer db.Close()

	// Get start of the day
	starOfTheDay := date.Bod(timestamp)

	// Copy candles from temp table into production
	log.Println("Delete candles data from: ", starOfTheDay)
	queryStr := `
	DELETE FROM candle
	WHERE ts > $1;`
	// Execute query
	_, err = db.Exec(queryStr, starOfTheDay)
	// Check err
	if err != nil {
		return err
	}

	return nil
}
