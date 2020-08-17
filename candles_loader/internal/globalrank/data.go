package globalrank

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
