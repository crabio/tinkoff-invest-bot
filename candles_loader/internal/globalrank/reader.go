package globalrank

import (
	"os"

	"github.com/gocarina/gocsv"
	"github.com/iakrevetkho/tinkoff-invest-bot/candles_loader/internal/data"
)

// ReadGlobalRankCsv reads global rank  data from CSV file
func ReadGlobalRankCsv(filePath string) (globalRanks []data.GlobalRank, err error) {
	// Open file
	in, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer in.Close()

	// Parse file
	if err := gocsv.UnmarshalFile(in, &globalRanks); err != nil {
		return nil, err
	}

	return
}
