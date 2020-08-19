package main

import (
	"flag"
	"github.com/iakrevetkho/tinkoff-invest-bot/candles_loader/internal/db"
	"log"
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

// Configuration from file Structure
var configurationFromFile = config.ConfigurationFile{}

// Configuration from environment variables Structure
var configurationFromEnv = config.ReadFromEnv()

func main() {
	// Setup logger
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)

	// Parse Arguments
	flag.Parse()

	log.Println("configurationFilePathPtr: ", configurationFilePathPtr)

	// Init error var
	var err error
	// Parse Config from JSON
	configurationFromFile, err = config.ReadFromFile(*configurationFilePathPtr)
	// Check error
	if err != nil {
		log.Println(err)
	}

	// Create database configuration
	dbConfiguration := db.Configuration{
		Type:     configurationFromEnv.DbType,
		User:     configurationFromEnv.DbUser,
		Password: configurationFromEnv.DbPassword,
		Hostname:  configurationFromEnv.DbHostname,
		Port:     configurationFromEnv.DbPort,
		DbName:   configurationFromEnv.DbName}

	// Wait DB init
	err = db.WaitDbInit(dbConfiguration, 10)
	// Check error
	if err != nil {
		log.Fatalln(err)
	}

	// Read Global Rank Companies rating
	globalRanks, err := globalrank.ReadGlobalRankCsv(configurationFromEnv.GlobalRankCsvFilePath)
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

	// Get latest loaded day from DB
	lattestTimestamp, err := db.GetLattestLoadedDay(dbConfiguration, configurationFromFile.StartLoadDate)
	// Check error
	if err != nil {
		log.Println("Lattest date in loaded days wasn't found")
		lattestTimestamp = configurationFromFile.StartLoadDate
	}

	// Delete first day data
	db.DeleteCandlesFromDay(dbConfiguration, lattestTimestamp)

	// Delete uploaded days
	db.DeleteUploadedDaysFromDay(dbConfiguration, lattestTimestamp)

	// Forever upload new data for day
	for true {
		// Generate date sequence
		daySequence := date.GenerateDaySequence(lattestTimestamp, time.Now())

		// Itterate over all days in sequence
		for dateX, date := range daySequence {
			log.Printf("Load data for date %d/%d: %v", dateX, len(daySequence), date)
			// Get candles for all instruments for one day
			for instrumentX, instrument := range globalRankInstuments {
				log.Printf("Load data for instrument %d/%d: '%v'", instrumentX, len(globalRankInstuments), instrument.Name)
				// Get candles
				candles, err := tinkoff.GetCandlesPerDay(configurationFromFile.ProductionToken,
					instrument,
					sdk.CandleInterval(configurationFromEnv.CandleInterval),
					date,
					configurationFromEnv.MaxAttempts)
				// Check error
				if err != nil {
					log.Println("Maybe problem with production token: ",
						configurationFromFile.ProductionToken,
						" or Internet onnection.")
					log.Fatalln(err)
				}

				// Load candles into DB
				err = db.UploadCandlesIntoDB(dbConfiguration, candles)
				// Check error
				if err != nil {
					log.Fatalln(err)
				}
			}

			// Add uploaded day into DB
			err = db.UploadNewLoadedDayIntoDB(dbConfiguration, date)
			// Check error
			if err != nil {
				log.Fatalln(err)
			}
		}

		// Get latest loaded day from DB
		lattestTimestamp, err = db.GetLattestLoadedDay(dbConfiguration, configurationFromFile.StartLoadDate)
		// Check error
		if err != nil {
			log.Println("Lattest date in loaded days wasn't found")
			lattestTimestamp = configurationFromFile.StartLoadDate
		}
		// Go to next day
		lattestTimestamp = lattestTimestamp.AddDate(0, 0, 1)

		// Wait next day
		time.Sleep(1 * time.Hour)
	}
}

// MatchGlobalRankOnstruments match list of Global Rank and Tinkoff Instruments
func MatchGlobalRankOnstruments(globalRanks []globalrank.GlobalRank) (globalRankInstruments []sdk.Instrument) {
	// Get all Tinkoff Markets
	instruments, err := tinkoff.GetAllMarkets(configurationFromFile.ProductionToken)
	if err != nil {
		log.Println("Maybe problem with production token: ",
			configurationFromFile.ProductionToken,
			" or Internet onnection.")
		log.Fatalln(err)
	}

	// Match global rank and Tinkoff instruments
	for _, globalRank := range globalRanks {
		for _, instrument := range instruments {
			// Match
			if smetrics.JaroWinkler(globalRank.Name, instrument.Name, 0.7, 4) > 0.9 {
				// Add to globalRankInstruments
				globalRankInstruments = append(globalRankInstruments, instrument)
			}
		}
	}

	log.Printf("Founded %d/%d Tinkoff Instruments in GLobal Rank", len(globalRankInstruments), len(globalRanks))
	return
}
