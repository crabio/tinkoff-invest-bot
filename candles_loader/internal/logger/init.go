package logger

import (
	log "github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
)

// Init logger settings
func Init() {
	// Log as default ASCII formatter.
	log.SetFormatter(
		&easy.Formatter{
			TimestampFormat: "2006-01-02 15:04:05",
			LogFormat:       "[%lvl%]: %time% - %msg%\n",
		})

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)
}
