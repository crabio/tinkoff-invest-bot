# Candles Loader

Programm for loading historical candles from Tinkoff API.

## Presequencies

* Golang > 1.12

## Configuration

Config file example:
```json
{
    "SandboxToken": "<YOUR TOKEN>",
    "ProductionToken": "YOUR TOKEN",
    "GlobalRankCsvFilePath": "../data/companies_rank.csv",
    "StartLoadDate": "2012-04-23T18:25:43.511Z"
}
```

## Development

For building use: `bash scripts/build.sh`

### Formatting

We are using `gofmt -w yourcode.go` for code formatting.

## TODOs

* App loads data from date from config
* When app is reloaded, it's find last date in table, go bacck for 1 day, delete newest day and start loading again