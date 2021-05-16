#!/bin/sh

docker-compose --file docker-compose.test.dependencies.yml up -d
docker-compose --file docker-compose.test.yml up
docker stop guideliner-db-test
