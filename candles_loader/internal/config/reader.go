package config

import (
	"github.com/tkanos/gonfig"
	"log"
	"os"
	"strconv"
)

// ReadFromFile Read Configuration struct from JSON file
func ReadFromFile(filePath string) (configuration ConfigurationFile, err error) {
	// Parse Config from JSON
	err = gonfig.GetConf(filePath, &configuration)

	// Print read config
	log.Println("From file ProductionToken = ", configuration.ProductionToken)
	log.Println("From file SandboxToken = ", configuration.SandboxToken)
	log.Println("From file StartLoadDate = ", configuration.StartLoadDate)

	return
}

// ReadFromEnv Read Configuration struct from environment variables
func ReadFromEnv() (configuration ConfigurationEnv) {
	// Read GLOBAL_RANK_CSV_FILE_PATH flag
	envGlobalRankCsvFilePath := os.Getenv("GLOBAL_RANK_CSV_FILE_PATH")
	// Check exists
	if envGlobalRankCsvFilePath != "" {
		log.Println("From ENV GlobalRankCsvFilePath = ", envGlobalRankCsvFilePath)
		configuration.GlobalRankCsvFilePath = envGlobalRankCsvFilePath
	}

	// Read MAX_ATTEMPTS flag
	envMaxAttemptsStr := os.Getenv("MAX_ATTEMPTS")
	// Check exists
	if envMaxAttemptsStr != "" {
		envMaxAttempts, err := strconv.Atoi(envMaxAttemptsStr)
		// Check err
		if err != nil {
			log.Println("Error: ", err)
		} else {
			log.Println("From ENV MaxAttempts = ", uint(envMaxAttempts))
			configuration.MaxAttempts = uint(envMaxAttempts)
		}
	}

	// Read CANDLE_INTERVAL flag
	envCandleInterval := os.Getenv("CANDLE_INTERVAL")
	// Check exists
	if envCandleInterval != "" {
		log.Println("From ENV DbType = ", envCandleInterval)
		configuration.CandleInterval = envCandleInterval
	}

	// Read DB_TYPE flag
	envDbType := os.Getenv("DB_TYPE")
	// Check exists
	if envDbType != "" {
		log.Println("From ENV DbType = ", envDbType)
		configuration.DbType = envDbType
	}

	// Read DB_USER flag
	envDbUser := os.Getenv("DB_USER")
	// Check exists
	if envDbUser != "" {
		log.Println("From ENV DbUser = ", envDbUser)
		configuration.DbUser = envDbUser
	}

	// Read DB_PASSWORD flag
	envDbPassword := os.Getenv("DB_PASSWORD")
	// Check exists
	if envDbPassword != "" {
		log.Println("From ENV DbPassword = ", envDbPassword)
		configuration.DbPassword = envDbPassword
	}

	// Read DB_HOSTNAME flag
	envDbHosname := os.Getenv("DB_HOSTNAME")
	// Check exists
	if envDbHosname != "" {
		log.Println("From ENV DbHosname = ", envDbHosname)
		configuration.DbHosname = envDbHosname
	}

	// Read DB_PORT flag
	envDbPortStr := os.Getenv("DB_PORT")
	// Check exists
	if envDbPortStr != "" {
		envDbPort, err := strconv.Atoi(envDbPortStr)
		// Check err
		if err != nil {
			log.Println("Error: ", err)
		} else {
			log.Println("From ENV DbPort = ", envDbPort)
			configuration.DbPort = uint(envDbPort)
		}
	}

	// Read DB_NAME flag
	envDbName := os.Getenv("DB_NAME")
	// Check exists
	if envDbName != "" {
		log.Println("From ENV DbName = ", envDbName)
		configuration.DbName = envDbName
	}

	return
}
