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
    "StartLoadDate": "2020-01-02T00:00:00.000Z",
    "MaxAttempts": 10,
    "DbType": "postgres",
    "DbUser": "postgres",
    "DbPassword": "postgres",
    "DbHosname": "timescaledb",
    "DbPort": 5432,
    "DbName": "tinkoff"
}
```

## Build

For building use: `bash scripts/build.sh`

## Docker

For run app in docker use: `bash scripts/build_docker.sh`

Application has ENV flags for configuring all required config fields:

* SANDBOX_TOKEN ""
* PRODUCTION_TOKEN ""
* GLOBAL_RANK_CSV_FILE_PATH "/data/companies_rank.csv"
* START_LOAD_DATE "2020-01-01T00:00:00.000Z"
* MAX_ATTEMPTS 10
* DB_TYPE "postgres"
* DB_USER "postgres"
* DB_PASSWORD "postgres"
* DB_HOSTNAME "timescaledb"
* DB_PORT 5432
* DB_NAME "tinkoff"

For running use: `docker run -v /data/tinkoff-invest-bot/data/:/data -e PRODUCTION_TOKEN="<your token>" candles-loader`

### Formatting

We are using `gofmt -w yourcode.go` for code formatting.

## TODOs