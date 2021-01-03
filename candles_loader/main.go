package main

import (
	"flag"
	"time"

	log "github.com/sirupsen/logrus"

	sdk "github.com/TinkoffCreditSystems/invest-openapi-go-sdk"
	"github.com/iakrevetkho/tinkoff-invest-bot/candles_loader/internal/config"
	"github.com/iakrevetkho/tinkoff-invest-bot/candles_loader/internal/date"
	"github.com/iakrevetkho/tinkoff-invest-bot/candles_loader/internal/db"
	"github.com/iakrevetkho/tinkoff-invest-bot/candles_loader/internal/globalrank"
	"github.com/iakrevetkho/tinkoff-invest-bot/candles_loader/internal/logger"
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
	logger.Init()

	log.Infoln("Start candles loader.")

	// Parse Arguments
	flag.Parse()

	log.Infoln("configurationFilePath: ", *configurationFilePathPtr)

	// Init error var
	var err error
	// Parse Config from JSON
	configurationFromFile, err = config.ReadFromFile(*configurationFilePathPtr)
	// Check error
	if err != nil {
		log.Fatal(err)
	}

	// Create database configuration
	dbConfiguration := db.Configuration{
		Type:     configurationFromEnv.DbType,
		User:     configurationFromEnv.DbUser,
		Password: configurationFromEnv.DbPassword,
		Hostname: configurationFromEnv.DbHostname,
		Port:     configurationFromEnv.DbPort,
		DbName:   configurationFromEnv.DbName}

	// Wait DB init
	err = db.WaitDbInit(dbConfiguration, 10)
	// Check error
	if err != nil {
		log.Fatal(err)
	}

	// Read Global Rank Companies rating
	globalRanks, err := globalrank.ReadGlobalRankCsv(configurationFromEnv.GlobalRankCsvFilePath)
	// Check error
	if err != nil {
		log.Fatal(err)
	}

	// Match global rank and Tinkoff instruments
	globalRankInstuments := MatchGlobalRankOnstruments(globalRanks)

	// Load instruments into DB
	err = db.UploadNewInstrumentsIntoDB(dbConfiguration, globalRankInstuments)
	// Check error
	if err != nil {
		log.Fatal(err)
	}

	// Get latest loaded day from DB
	lattestTimestamp, err := db.GetLattestLoadedDay(dbConfiguration, configurationFromFile.StartLoadDate)
	// Check error
	if err != nil {
		log.Debug("Lattest date in loaded days wasn't found")
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
			log.Debugf("Load data for date %d/%d: %v", dateX, len(daySequence), date)
			// Get candles for all instruments for one day
			for instrumentX, instrument := range globalRankInstuments {
				log.Tracef("Load data for instrument %d/%d: '%v'", instrumentX, len(globalRankInstuments), instrument.Name)
				// Get candles
				candles, err := tinkoff.GetCandlesPerDay(configurationFromFile.ProductionToken,
					instrument,
					sdk.CandleInterval(configurationFromEnv.CandleInterval),
					date,
					configurationFromEnv.MaxAttempts)
				// Check error
				if err != nil {
					log.Fatal("Maybe problem with production token: ",
						configurationFromFile.ProductionToken,
						" or Internet onnection.")
				}

				// Load candles into DB
				err = db.UploadCandlesIntoDB(dbConfiguration, candles)
				// Check error
				if err != nil {
					log.Fatal(err)
				}
			}

			// Add uploaded day into DB
			err = db.UploadNewLoadedDayIntoDB(dbConfiguration, date)
			// Check error
			if err != nil {
				log.Fatal(err)
			}
		}

		// Get latest loaded day from DB
		lattestTimestamp, err = db.GetLattestLoadedDay(dbConfiguration, configurationFromFile.StartLoadDate)
		// Check error
		if err != nil {
			log.Debug("Lattest date in loaded days wasn't found")
			lattestTimestamp = configurationFromFile.StartLoadDate
		}
		// Go to next day
		lattestTimestamp = lattestTimestamp.AddDate(0, 0, 1)

		// Wait next day
		time.Sleep(1 * time.Hour)
	}

	log.Infoln("Stop candles loader.")
}

// MatchGlobalRankOnstruments match list of Global Rank and Tinkoff Instruments
func MatchGlobalRankOnstruments(globalRanks []globalrank.GlobalRank) (globalRankInstruments []sdk.Instrument) {
	// Get all Tinkoff Markets
	instruments, err := tinkoff.GetAllMarkets(configurationFromFile.ProductionToken)
	if err != nil {
		log.Fatal("Maybe problem with production token: ",
			configurationFromFile.ProductionToken,
			" or Internet onnection.")
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

	log.Debugf("Founded %d/%d Tinkoff Instruments in GLobal Rank", len(globalRankInstruments), len(globalRanks))
	return
}
