# Docker

## Prepare

For stopping all containers execute: `docker stop $(docker ps -a -q)`

For removing all containers and data execute: `docker system prune --volumes`

## Configuration

Create file `.env` with reqired TOKENS and optional another settings or add all settings into `docker-compose.yml` file.

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
