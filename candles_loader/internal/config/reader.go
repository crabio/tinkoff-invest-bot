package config

import (
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/tkanos/gonfig"
)

// ReadFromFile Read Configuration struct from JSON file
func ReadFromFile(filePath string) (configuration ConfigurationFile, err error) {
	// Parse Config from JSON
	err = gonfig.GetConf(filePath, &configuration)

	// Print read config
	log.Debugln("From file ProductionToken = ", configuration.ProductionToken)
	log.Debugln("From file SandboxToken = ", configuration.SandboxToken)
	log.Debugln("From file StartLoadDate = ", configuration.StartLoadDate)

	return
}

// ReadFromEnv Read Configuration struct from environment variables
func ReadFromEnv() (configuration ConfigurationEnv) {
	// Read GLOBAL_RANK_CSV_FILE_PATH flag
	envGlobalRankCsvFilePath := os.Getenv("GLOBAL_RANK_CSV_FILE_PATH")
	// Check exists
	if envGlobalRankCsvFilePath != "" {
		log.Debugln("From ENV GlobalRankCsvFilePath = ", envGlobalRankCsvFilePath)
		configuration.GlobalRankCsvFilePath = envGlobalRankCsvFilePath
	}

	// Read MAX_ATTEMPTS flag
	envMaxAttemptsStr := os.Getenv("MAX_ATTEMPTS")
	// Check exists
	if envMaxAttemptsStr != "" {
		envMaxAttempts, err := strconv.Atoi(envMaxAttemptsStr)
		// Check err
		if err != nil {
			log.Debugln("Error: ", err)
		} else {
			log.Debugln("From ENV MaxAttempts = ", uint(envMaxAttempts))
			configuration.MaxAttempts = uint(envMaxAttempts)
		}
	}

	// Read CANDLE_INTERVAL flag
	envCandleInterval := os.Getenv("CANDLE_INTERVAL")
	// Check exists
	if envCandleInterval != "" {
		log.Debugln("From ENV DbType = ", envCandleInterval)
		configuration.CandleInterval = envCandleInterval
	}

	// Read DB_TYPE flag
	envDbType := os.Getenv("DB_TYPE")
	// Check exists
	if envDbType != "" {
		log.Debugln("From ENV DbType = ", envDbType)
		configuration.DbType = envDbType
	}

	// Read DB_USER flag
	envDbUser := os.Getenv("DB_USER")
	// Check exists
	if envDbUser != "" {
		log.Debugln("From ENV DbUser = ", envDbUser)
		configuration.DbUser = envDbUser
	}

	// Read DB_PASSWORD flag
	envDbPassword := os.Getenv("DB_PASSWORD")
	// Check exists
	if envDbPassword != "" {
		log.Debugln("From ENV DbPassword = ", envDbPassword)
		configuration.DbPassword = envDbPassword
	}

	// Read DB_HOSTNAME flag
	envDbHostname := os.Getenv("DB_HOSTNAME")
	// Check exists
	if envDbHostname != "" {
		log.Debugln("From ENV DbHostname = ", envDbHostname)
		configuration.DbHostname = envDbHostname
	}

	// Read DB_PORT flag
	envDbPortStr := os.Getenv("DB_PORT")
	// Check exists
	if envDbPortStr != "" {
		envDbPort, err := strconv.Atoi(envDbPortStr)
		// Check err
		if err != nil {
			log.Debugln("Error: ", err)
		} else {
			log.Debugln("From ENV DbPort = ", envDbPort)
			configuration.DbPort = uint(envDbPort)
		}
	}

	// Read DB_NAME flag
	envDbName := os.Getenv("DB_NAME")
	// Check exists
	if envDbName != "" {
		log.Debugln("From ENV DbName = ", envDbName)
		configuration.DbName = envDbName
	}

	return
}
