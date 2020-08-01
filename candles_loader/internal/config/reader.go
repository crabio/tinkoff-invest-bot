package config

import (
	"log"

	"github.com/tkanos/gonfig"
)

// ReadFromFile Read Configuration struct from JSON file
func ReadFromFile(filePath string) (configuration Configuration) {
	// Parse Config from JSON
	err := gonfig.GetConf(filePath, &configuration)
	// Check error
	if err != nil {
		log.Fatalln(err)
	}
	return
}
