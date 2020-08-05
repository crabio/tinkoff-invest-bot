# Candles Loader

Programm for loading historical candles from Tinkoff API.

## Presequencies

* Golang > 1.12

## Development

For building use: `bash scripts/build.sh`

### Formatting

We are using `gofmt -w yourcode.go` for code formatting.

## TODOs

* App loads data from date from config
* When app is reloaded, it's find last date in table, go bacck for 1 day, delete newest day and start loading again