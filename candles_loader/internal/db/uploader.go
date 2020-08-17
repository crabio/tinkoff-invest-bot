package db

import (
	"fmt"
	"github.com/TinkoffCreditSystems/invest-openapi-go-sdk"
	"github.com/lib/pq"
	"log"
	"time"
)

// UploadNewInstrumentsIntoDB upload instruments meta information into data base
func UploadNewInstrumentsIntoDB(config Configuration, instruments []sdk.Instrument) (err error) {
	// Connect to DB
	db, err := CreateDbConnection(config)
	// Check err
	if err != nil {
		return err
	}
	// At the end close connetion
	defer db.Close()

	// Create Temp table
	queryStr := `CREATE TEMPORARY TABLE temp_instrument(
		id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
		figi VARCHAR(255) NOT NULL,
		ticker VARCHAR(255) NOT NULL,
		name VARCHAR(255) NOT NULL,
		min_price_increment FLOAT NOT NULL,
		currency VARCHAR(255) NOT NULL,
		type VARCHAR(255) NOT NULL,
		global_rank BOOL NOT NULL
	);`
	// Execute query
	_, err = db.Exec(queryStr)
	// Check err
	if err != nil {
		return err
	}

	// Start transaction
	transaction, err := db.Begin()
	// Check err
	if err != nil {
		return err
	}

	// Prepare copy query into temp table
	statement, err := transaction.Prepare(pq.CopyIn(
		"temp_instrument", "figi", "ticker", "name", "min_price_increment", "currency", "type", "global_rank"))
	if err != nil {
		log.Fatal(err)
	}

	// Insert each instrument
	for _, instrument := range instruments {
		// Execute query
		_, err = statement.Exec(instrument.FIGI,
			instrument.Ticker,
			instrument.Name,
			instrument.MinPriceIncrement,
			instrument.Currency,
			instrument.Type,
			true)
		// Check error
		if err != nil {
			return err
		}
	}
	// Final exec to close loading process
	_, err = statement.Exec()
	// Check error
	if err != nil {
		return err
	}

	// Close statement
	err = statement.Close()
	// Check error
	if err != nil {
		return err
	}

	// Commit transaction
	err = transaction.Commit()
	// Check error
	if err != nil {
		return err
	}

	// Copy only new rows from temp table into production
	queryStr = `
	INSERT INTO instrument (figi, ticker, name, min_price_increment, currency, type, global_rank)
	SELECT temp.figi, temp.ticker, temp.name, temp.min_price_increment, temp.currency, temp.type, temp.global_rank
	FROM temp_instrument as temp
	LEFT JOIN instrument as instrument
	ON temp.figi = instrument.figi
	WHERE instrument.id IS NULL;`
	// Execute query
	_, err = db.Exec(queryStr)
	// Check err
	if err != nil {
		return err
	}

	return nil
}

// UploadNewLoadedDayIntoDB upload loaded day into data base list
func UploadNewLoadedDayIntoDB(config Configuration, loadedDay time.Time) (err error) {
	// Connect to DB
	db, err := CreateDbConnection(config)
	// Check err
	if err != nil {
		return err
	}
	// At the end close connetion
	defer db.Close()

	// Start transaction
	transaction, err := db.Begin()
	// Check err
	if err != nil {
		return err
	}

	// Prepare copy query
	statement, err := transaction.Prepare(pq.CopyIn(
		"candle_loaded_day", "day"))
	if err != nil {
		log.Fatal(err)
	}

	// Insert day
	_, err = statement.Exec(loadedDay)
	// Check error
	if err != nil {
		return err
	}

	// Close statement
	err = statement.Close()
	// Check error
	if err != nil {
		return err
	}

	// Commit transaction
	err = transaction.Commit()
	// Check error
	if err != nil {
		return err
	}

	return
}

// UploadCandlesIntoDB upload candles into data base
func UploadCandlesIntoDB(config Configuration, candles []sdk.Candle) (err error) {
	// Connect to DB
	db, err := CreateDbConnection(config)
	// Check err
	if err != nil {
		return err
	}
	// At the end close connetion
	defer db.Close()

	// Create Temp table
	queryStr := `CREATE TEMPORARY TABLE temp_candle(
		ts timestamptz NOT NULL,
		instrument_figi VARCHAR(255) NOT NULL,
		interval_name VARCHAR(255) NOT NULL,
		open_price REAL NULL,
		close_price REAL NULL,
		high_price REAL NULL,
		low_price REAL NULL,
		volume REAL NULL
	);`
	// Execute query
	_, err = db.Exec(queryStr)
	// Check err
	if err != nil {
		return err
	}

	// Start transaction
	transaction, err := db.Begin()
	// Check err
	if err != nil {
		return err
	}

	// Prepare copy query into temp table
	statement, err := transaction.Prepare(pq.CopyIn(
		"temp_candle", "ts", "instrument_figi", "interval_name",
		"open_price", "close_price", "high_price", "low_price", "volume"))
	if err != nil {
		log.Fatal(err)
	}

	// Insert each instrument
	for _, candle := range candles {
		// Execute query
		_, err = statement.Exec(candle.TS,
			candle.FIGI,
			candle.Interval,
			candle.OpenPrice,
			candle.ClosePrice,
			candle.HighPrice,
			candle.LowPrice,
			candle.Volume)
		// Check error
		if err != nil {
			return err
		}
	}
	// Final exec to close loading process
	_, err = statement.Exec()
	// Check error
	if err != nil {
		return err
	}

	// Close statement
	err = statement.Close()
	// Check error
	if err != nil {
		return err
	}

	// Commit transaction
	err = transaction.Commit()
	// Check error
	if err != nil {
		return err
	}

	// Add new candles intervals if found
	queryStr = `
	INSERT INTO candle_interval (name)
	SELECT DISTINCT temp.interval_name
	FROM temp_candle as temp
	LEFT JOIN candle_interval as interval
	ON temp.interval_name = interval.name
	WHERE interval.id IS NULL;`
	// Execute query
	_, err = db.Exec(queryStr)
	// Check err
	if err != nil {
		return err
	}

	// Get unknown Instruments
	queryStr = `
	SELECT temp.interval_name
	FROM temp_candle as temp
	LEFT JOIN instrument as instrument
	ON temp.instrument_figi = instrument.figi
	WHERE instrument.id IS NULL;`
	// Execute query
	rows, err := db.Query(queryStr)
	// Check err
	if err != nil {
		return err
	}
	// Create list for unknown instrument names
	unknownInstrumentNames := []string{}
	// Get names from query
	for rows.Next() {
		// Init buffer
		var unknownInstrumentName string
		// Convert data from query itterator to buffer
		err := rows.Scan(&unknownInstrumentName)
		// Check error
		if err != nil {
			return err
		}
		// Append new name
		unknownInstrumentNames = append(unknownInstrumentNames, unknownInstrumentName)
	}
	// Check rows count
	if len(unknownInstrumentNames) != 0 {
		return error(fmt.Errorf("Not all instruments are found in DB: %v+", unknownInstrumentNames))
	}
	// Close query itterator
	err = rows.Close()
	// Check err
	if err != nil {
		return err
	}

	// Copy candles from temp table into production
	queryStr = `
	INSERT INTO candle (ts, instrument_id, interval_id,
		open_price, close_price, high_price, low_price, volume)
	SELECT temp.ts, instrument.id, interval.id,
		temp.open_price, temp.close_price, temp.high_price, temp.low_price, temp.volume
	FROM temp_candle as temp
	LEFT JOIN instrument as instrument ON temp.instrument_figi = instrument.figi
	LEFT JOIN candle_interval as interval ON temp.interval_name = interval.name;`
	// Execute query
	_, err = db.Exec(queryStr)
	// Check err
	if err != nil {
		return err
	}

	return nil
}
