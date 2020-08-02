package main

import (
	"flag"
	"log"

	"github.com/iakrevetkho/tinkoff-invest-bot/candles_loader/internal/config"
	"github.com/iakrevetkho/tinkoff-invest-bot/candles_loader/internal/globalrank"
	"github.com/iakrevetkho/tinkoff-invest-bot/candles_loader/internal/tinkoff"
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
	_, err := globalrank.ReadGlobalRankCsv(configuration.GlobalRankCsvFilePath)
	if err != nil {
		log.Fatalln(err)
	}

	// Get all Tinkoff Markets
	instrumentsNameMap := tinkoff.GetAllMarketsMap(configuration.ProductionToken)

	log.Println(instrumentsNameMap)

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
