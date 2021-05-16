#!/bin/sh

SCRIPTPATH="$( cd -- "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"

go mod vendor

make guideliner-build
make migrations-build
make clean-db-build

/bin/sh $SCRIPTPATH/../sh/wait_for_postgresql.sh

echo >&2 "Clean db..."
"$SCRIPTPATH"/../../bin/clean_db

echo >&2 "Run migrations..."
"$SCRIPTPATH"/../../bin/migrations

echo >&2 "Run server..."
"$SCRIPTPATH"/../../bin/guideliner &

/bin/sh "$SCRIPTPATH"/../sh/wait_for_server.sh

make tests-integration

