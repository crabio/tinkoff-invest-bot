package db

import (
	"database/sql"
	"fmt"
	"github.com/TinkoffCreditSystems/invest-openapi-go-sdk"
	"github.com/lib/pq"
	"log"
)

// UploadInstrumentsIntoDB upload instruments meta information into data base
func UploadInstrumentsIntoDB(config Configuration, instruments []sdk.Instrument) (err error) {
	// Create onnection string
	connectionString := fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=disable",
		config.DbType, config.User, config.Password, config.Hosname, config.Port, config.DbName)

	// Connect to DB
	db, err := sql.Open("db", connectionString)
	// Check err
	if err != nil {
		return err
	}

	// Create Temp table
	queryStr := `CREATE TEMPORARY TABLE temp_instrument(
		id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
		figi VARCHAR(255) NOT NULL,
		ticker VARCHAR(255) NOT NULL,
		name VARCHAR(255) NOT NULL,
		min_price_increment FLOAT NOT NULL,
		currency VARCHAR(10) NOT NULL,
		type VARCHAR(50) NOT NULL
	);`
	// Execute query
	_, err = db.Exec(queryStr)
	// Check err
	if err != nil {
		return err
	}

	// Start transaction
	txn, err := db.Begin()
	// Check err
	if err != nil {
		return err
	}

	// Prepare copy query into temp table
	stmt, err := txn.Prepare(pq.CopyIn(
		"temp_instrument", "figi", "ticker", "name", "min_price_increment", "currency", "type"))
	if err != nil {
		log.Fatal(err)
	}

	// Insert each instrument
	for _, instrument := range instruments {
		// Execute query
		_, err = stmt.Exec(instrument.FIGI,
			instrument.Ticker,
			instrument.Name,
			instrument.MinPriceIncrement,
			instrument.Currency,
			instrument.Type)
		// Check error
		if err != nil {
			return err
		}
	}
	// Final exec to close loading process
	_, err = stmt.Exec()
	// Check error
	if err != nil {
		return err
	}

	// Close statement
	err = stmt.Close()
	// Check error
	if err != nil {
		return err
	}

	// Commit transaction
	err = txn.Commit()
	// Check error
	if err != nil {
		return err
	}

	// Copy only new rows from temp table into production
	queryStr = `
	INSERT INTO $1 (figi, ticker, name, min_price_increment, currency, type)
	SELECT temp.figi, temp.ticker, temp.name, temp.min_price_increment, temp.currency, temp.type
	FROM temp_instrument as temp
	LEFT JOIN $1 as prod
	ON temp.figi = prod.figi
	WHERE prod.id IS NULL;`
	// Execute query
	_, err = db.Exec(queryStr, config.InstrumentsTableName)
	// Check err
	if err != nil {
		return err
	}

	// Close connection
	err = db.Close()
	// Check error
	if err != nil {
		return err
	}

	return nil
}
