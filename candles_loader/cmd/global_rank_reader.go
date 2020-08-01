package main

import (
	"os"

	"github.com/gocarina/gocsv"
)

// GlobalRank - information about company rank in global Forbes rating
type GlobalRank struct {
	Name        string  `csv:"Company"`
	Country     string  `csv:"Country"`
	Industry    string  `csv:"Industry"`
	Sales       float32 `csv:"Sales"`
	Profits     float32 `csv:"Profits"`
	Assets      float32 `csv:"Assets"`
	MarketValue float32 `csv:"Market Value"`
}

func readGlobalRankParquetr(filePath string) (globalRanks []GlobalRank, err error) {
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
