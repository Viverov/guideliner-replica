#!/bin/sh
docker-compose --file docker-compose.test.dependencies.yml up -d
docker-compose --file docker-compose.test.yml up --exit-code-from server
status=$?
docker stop guideliner-db-test
exit $status
