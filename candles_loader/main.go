package main

import (
	"flag"
	"log"

	"github.com/iakrevetkho/tinkoff-invest-bot/candles_loader/internal/config"
	"github.com/iakrevetkho/tinkoff-invest-bot/candles_loader/internal/globalrank"
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

	globalRanks, err := globalrank.ReadGlobalRankCsv(configuration.GlobalRankCsvFilePath)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(globalRanks)

	// if configuration.IsSandbox {
	// 	sandboxRest()
	// } else {
	// 	rest()
	// }
}
