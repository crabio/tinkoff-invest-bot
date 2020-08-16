# Docker

## Prepare

For stopping all containers execute: `docker stop $(docker ps -a -q)`

For removing all containers and data execute: `docker system prune --volumes`

## Configuration

### Overal

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

### Superset

For Superset (BI) configuration we use `superset.env`.

## Run

For running, execute:

* `docker-compose up -d`

## Stop

For stopping and delete DB, execute:

* `docker-compose down --remove-orphans`

## Monitoring

For monitoring we are using `dockprom`.
For running, execute:

* `cd dockprom`
* `ADMIN_USER=admin ADMIN_PASSWORD=admin docker-compose up -d`
