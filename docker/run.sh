#!/bin/bash

echo 'Deploy app.'

docker-compose up -d

echo 'Deploy monitoring.'

( cd dockprom && ADMIN_USER=admin ADMIN_PASSWORD=admin docker-compose up -d )
