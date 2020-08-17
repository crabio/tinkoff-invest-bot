# Candles Loader

Programm for loading historical candles from Tinkoff API.

## Structure

* main.go - Main Golang source file
* go.mod, go.sum - Golang dependenccies files
* Dockerfile - fFFIle for dockerize application
* build - Folder for storing built binaries
* internal - Private sources/libraries

## Presequencies

* Golang > 1.12

## Configuration

Config file example:

```json
{
    "SandboxToken": "<YOUR TOKEN>",
    "ProductionToken": "YOUR TOKEN",
    "StartLoadDate": "2020-01-02T00:00:00.000Z"
}
```

Configuration  file  should be placed in `/tinkoff-invest-bot/config/config.json`
as signle configuration point for all services.

## Build

For building use: `bash scripts/build.sh`

## Docker

For run app in docker use: `bash scripts/build_docker.sh`

Application has ENV flags for configuring all required config fields:

* GLOBAL_RANK_CSV_FILE_PATH "/data/companies_rank.csv"
* MAX_ATTEMPTS 10
* CANDLE_INTERVAL "15min" - Allowed candles interval are: "1min","2min","3min","5min","10min","15min","30min","hour","2hour","4hour","day","week","month"
* DB_TYPE "postgres"
* DB_USER "postgres"
* DB_PASSWORD "postgres"
* DB_HOSTNAME "timescaledb"
* DB_PORT 5432
* DB_NAME "tinkoff"

For running use: `docker run -v /tinkoff-invest-bot/data/:/data --env-file .env --name candles-loader -d candles-loader`

### Formatting

We are using `gofmt -w yourcode.go` for code formatting.

## TODOs
