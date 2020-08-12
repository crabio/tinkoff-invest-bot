package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

// ExtendWithEnvFlags extends config in arguments with Env flags and return new config
func ExtendWithEnvFlags(oldConfiguration Configuration) (newConfiguration Configuration) {
	// Copy configuration
	newConfiguration = oldConfiguration

	// Read SANDBOX_TOKEN flag
	envSandboxToken := os.Getenv("SANDBOX_TOKEN")
	// Check exists
	if envSandboxToken != "" {
		log.Println("From ENV SandboxToken = ", envSandboxToken)
		newConfiguration.SandboxToken = envSandboxToken
	}

	// Read PRODUCTION_TOKEN flag
	envProductionToken := os.Getenv("PRODUCTION_TOKEN")
	// Check exists
	if envProductionToken != "" {
		log.Println("From ENV ProductionToken = ", envProductionToken)
		newConfiguration.ProductionToken = envProductionToken
	}

	// Read GLOBAL_RANK_CSV_FILE_PATH flag
	envGlobalRankCsvFilePath := os.Getenv("GLOBAL_RANK_CSV_FILE_PATH")
	// Check exists
	if envGlobalRankCsvFilePath != "" {
		log.Println("From ENV GlobalRankCsvFilePath = ", envGlobalRankCsvFilePath)
		newConfiguration.GlobalRankCsvFilePath = envGlobalRankCsvFilePath
	}

	// Read START_LOAD_DATE flag
	envStartLoadDateStr := os.Getenv("START_LOAD_DATE")
	// Check exists
	if envStartLoadDateStr != "" {
		envStartLoadDate, err := time.Parse(time.RFC3339, envStartLoadDateStr)
		// Check err
		if err != nil {
			log.Println("Error: ", err)
		} else {
			log.Println("From ENV StartLoadDate = ", envStartLoadDate)
			newConfiguration.StartLoadDate = envStartLoadDate
		}
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
			newConfiguration.MaxAttempts = uint(envMaxAttempts)
		}
	}

	// Read CANDLE_INTERVAL flag
	envCandleInterval := os.Getenv("CANDLE_INTERVAL")
	// Check exists
	if envCandleInterval != "" {
		log.Println("From ENV DbType = ", envCandleInterval)
		newConfiguration.CandleInterval = envCandleInterval
	}

	// Read DB_TYPE flag
	envDbType := os.Getenv("DB_TYPE")
	// Check exists
	if envDbType != "" {
		log.Println("From ENV DbType = ", envDbType)
		newConfiguration.DbType = envDbType
	}

	// Read DB_USER flag
	envDbUser := os.Getenv("DB_USER")
	// Check exists
	if envDbUser != "" {
		log.Println("From ENV DbUser = ", envDbUser)
		newConfiguration.DbUser = envDbUser
	}

	// Read DB_PASSWORD flag
	envDbPassword := os.Getenv("DB_PASSWORD")
	// Check exists
	if envDbPassword != "" {
		log.Println("From ENV DbPassword = ", envDbPassword)
		newConfiguration.DbPassword = envDbPassword
	}

	// Read DB_HOSTNAME flag
	envDbHosname := os.Getenv("DB_HOSTNAME")
	// Check exists
	if envDbHosname != "" {
		log.Println("From ENV DbHosname = ", envDbHosname)
		newConfiguration.DbHosname = envDbHosname
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
			newConfiguration.DbPort = uint(envDbPort)
		}
	}

	// Read DB_NAME flag
	envDbName := os.Getenv("DB_NAME")
	// Check exists
	if envDbName != "" {
		log.Println("From ENV DbName = ", envDbName)
		newConfiguration.DbName = envDbName
	}

	return
}
