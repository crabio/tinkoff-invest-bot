package config

import (
	"github.com/tkanos/gonfig"
)

// ReadFromFile Read Configuration struct from JSON file
func ReadFromFile(filePath string) (configuration Configuration, err error) {
	// Parse Config from JSON
	err = gonfig.GetConf(filePath, &configuration)

	return
}
