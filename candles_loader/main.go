package main

import (
	"flag"
	"github.com/iakrevetkho/tinkoff-invest-bot/candles_loader/internal/db"
	"log"
	"sync"
	"time"

	sdk "github.com/TinkoffCreditSystems/invest-openapi-go-sdk"
	"github.com/iakrevetkho/tinkoff-invest-bot/candles_loader/internal/config"
	"github.com/iakrevetkho/tinkoff-invest-bot/candles_loader/internal/date"
	"github.com/iakrevetkho/tinkoff-invest-bot/candles_loader/internal/globalrank"
	"github.com/iakrevetkho/tinkoff-invest-bot/candles_loader/internal/tinkoff"
	"github.com/xrash/smetrics"
)

// Configuration File Path. Default "config.json"
var configurationFilePathPtr = flag.String("c", "config.json", "Configuration File Path")

// Configuration Structure
var configuration = config.Configuration{}

// Database configuration
var dbConfiguration = db.Configuration{
	DbType:   "postgres",
	User:     "postgres",
	Password: "postgres",
	Hosname:  "localhost",
	Port:     5432,
	DbName:   "tinkoff"}

func main() {
	// Parse Arguments
	flag.Parse()

	// Parse Config from JSON
	configuration = config.ReadFromFile(*configurationFilePathPtr)

	// Read Global Rank Companies rating
	globalRanks, err := globalrank.ReadGlobalRankCsv(configuration.GlobalRankCsvFilePath)
	// Check error
	if err != nil {
		log.Fatalln(err)
	}

	// Match global rank and Tinkoff instruments
	globalRankInstuments := MatchGlobalRankOnstruments(globalRanks)

	// Load instruments into DB
	err = db.UploadNewInstrumentsIntoDB(dbConfiguration, globalRankInstuments)
	// Check error
	if err != nil {
		log.Fatalln(err)
	}

	// Get latest candle date in DB
	lattestTimestamp, err := db.GetLattestCandleTimestamp(dbConfiguration, configuration.StartLoadDate)
	// Check error
	if err != nil {
		log.Println("Lattest timestamp in candles wasn't found")
		lattestTimestamp = configuration.StartLoadDate
	}
	// Get start of the day
	lattestTimestamp = date.Bod(lattestTimestamp)

	// Delete first day data
	db.DeleteCandlesFromDay(dbConfiguration, lattestTimestamp)

	// Generate date sequence
	daySequence := date.GenerateDaySequence(lattestTimestamp, time.Now())

	// Itterate over all days in sequence
	for dateX, date := range daySequence {
		log.Printf("Load data for date %d/%d: %v", dateX, len(daySequence), date)
		// Get candles for all instruments for one day
		for instrumentX, instrument := range globalRankInstuments {
			log.Printf("Load data for instrument %d/%d: '%v'", instrumentX, len(globalRankInstuments), instrument.Name)
			// Get candles
			candles := tinkoff.GetCandlesPerDay(configuration.ProductionToken,
				instrument,
				sdk.CandleInterval15Min,
				date,
				configuration.MaxAttempts)

			// Load candles into DB
			err = db.UploadCandlesIntoDB(dbConfiguration, candles)
			// Check error
			if err != nil {
				log.Fatalln(err)
			}
		}
	}

	// if configuration.IsSandbox {
	// 	sandboxRest()
	// } else {
	// 	rest()
	// }
}

// MatchGlobalRankOnstruments match list of Global Rank and Tinkoff Instruments
func MatchGlobalRankOnstruments(globalRanks []globalrank.GlobalRank) (globalRankInstruments []sdk.Instrument) {
	// Get all Tinkoff Markets
	instruments, err := tinkoff.GetAllMarkets(configuration.ProductionToken)
	if err != nil {
		log.Fatalln(err)
	}

	// Match global rank and Tinkoff instruments
	counter := 0
	var wg sync.WaitGroup

	for _, globalRank := range globalRanks {
		// These goroutines share memory, but only for reading.
		wg.Add(1)

		go func(globalRank globalrank.GlobalRank) {

			for _, instrument := range instruments {
				// Match
				if smetrics.JaroWinkler(globalRank.Name, instrument.Name, 0.7, 4) > 0.9 {
					counter++
					// Add to globalRankInstruments
					globalRankInstruments = append(globalRankInstruments, instrument)
				}
			}

			wg.Done()
		}(globalRank)
	}
	wg.Wait()
	log.Printf("Founded %d/%d Tinkoff Instruments in GLobal Rank", counter, len(globalRanks))
	return
}
