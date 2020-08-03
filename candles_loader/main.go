package main

import (
	"flag"
	"log"
	"sync"

	sdk "github.com/TinkoffCreditSystems/invest-openapi-go-sdk"
	"github.com/iakrevetkho/tinkoff-invest-bot/candles_loader/internal/config"
	"github.com/iakrevetkho/tinkoff-invest-bot/candles_loader/internal/globalrank"
	"github.com/iakrevetkho/tinkoff-invest-bot/candles_loader/internal/tinkoff"
	"github.com/xrash/smetrics"
)

// Configuration File Path. Default "config.json"
var configurationFilePathPtr = flag.String("c", "config.json", "Configuration File Path")

// Configuration Structure
var configuration = config.Configuration{}

func main() {
	// Parse Arguments
	flag.Parse()

	// Parse Config from JSON
	configuration = config.ReadFromFile(*configurationFilePathPtr)

	// Read Global Rank Companies rating
	globalRanks, err := globalrank.ReadGlobalRankCsv(configuration.GlobalRankCsvFilePath)
	if err != nil {
		log.Fatalln(err)
	}

	// Get all Tinkoff Markets
	instruments, err := tinkoff.GetAllMarkets(configuration.ProductionToken)
	if err != nil {
		log.Fatalln(err)
	}

	// Match global rank and Tinkoff instruments
	counter := 0
	var globalRankInstruments []sdk.Instrument
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
	log.Println(globalRankInstruments)
	log.Printf("Founded %d/%d", counter, len(globalRanks))

	// // Get Tickets FIGI from Tinkoff
	// for _, globalRank := range globalRanks {
	// 	tinkoff.GetFigiByTicket(configuration.ProductionToken, globalRank.Name)
	// }

	// if configuration.IsSandbox {
	// 	sandboxRest()
	// } else {
	// 	rest()
	// }
}
