package tinkoff

import (
	"context"
	"log"
	"time"

	sdk "github.com/TinkoffCreditSystems/invest-openapi-go-sdk"
)

// GetAllMarkets creates list of all Tinkoff markets
func GetAllMarkets(token string) (instruments []sdk.Instrument, err error) {
	// Create REST Client
	client := sdk.NewRestClient(token)

	// Create context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("Get currency instruments")
	// Example: USD000UTSTOM - USD, EUR_RUB__TOM - EUR
	currencies, err := client.Currencies(ctx)
	if err != nil {
		return nil, err
	}
	// Add currencies
	instruments = append(instruments, currencies...)

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("Fet etf instruments")
	// Example: FXMM - Казначейские облигации США, FXGD - золото
	etfs, err := client.ETFs(ctx)
	if err != nil {
		return nil, err
	}
	// Add etfs
	instruments = append(instruments, etfs...)

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("Get bond instruments")
	// Example: SU24019RMFS0 - ОФЗ 24019
	bonds, err := client.Bonds(ctx)
	if err != nil {
		return nil, err
	}
	// Add bonds
	instruments = append(instruments, bonds...)

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("Add stock instruments")
	// Example: SBUX - Starbucks Corporation
	stocks, err := client.Stocks(ctx)
	if err != nil {
		return nil, err
	}
	// Add stocks
	instruments = append(instruments, stocks...)

	return
}

// GetCandlesPerDay receives candles per day with custom inverval for Tinkoff instrument
func GetCandlesPerDay(
	token string, instrument sdk.Instrument,
	interval sdk.CandleInterval,
	date time.Time,
	maxAttempts uint) (
	candles []sdk.Candle) {
	// Create REST Client
	client := sdk.NewRestClient(token)

	// Until not success ot attemps out of range
	var success bool = false
	var attempt uint = 0
	for !success {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Get candles for specific day
		var err error
		candles, err = client.Candles(ctx, date, date.AddDate(0, 0, 1), interval, instrument.FIGI)
		if err != nil {
			// Check attempts counter
			if attempt == maxAttempts {
				log.Fatalln("maxAttempts count reached")
			}
			// New attempt
			time.Sleep(1 * time.Second)
			attempt++
		} else {
			// Success request
			success = true
		}
	}

	return
}
