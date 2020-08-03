# Docker

## Prepare

For stopping all containers execute: `docker stop $(docker ps -a -q)`

For removing all containers and data execute: `docker system prune --volumes`

## Run

For running, execute:

* `docker-compose up -d`

## Monitoring

For monitoring we are using `dockprom`.
For running, execute:

* `cd dockprom`
* `ADMIN_USER=admin ADMIN_PASSWORD=admin docker-compose up -d`
