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

	// Start transaction
	txn, err := db.Begin()
	// Check err
	if err != nil {
		return err
	}

	// Prepare copy query into table
	stmt, err := txn.Prepare(pq.CopyIn(config.InstrumentsTableName,
		"figi", "ticker", "name", "min_price_increment", "currency", "type"))
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

	return nil
}
