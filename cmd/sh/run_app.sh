#!/bin/sh

SCRIPTPATH="$( cd -- "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"

/bin/sh $SCRIPTPATH/../sh/wait_for_postgresql.sh
status=$?
if [ $status -ne 0 ]; then
  exit $status
fi

echo >&2 "Run migrations..."
"$SCRIPTPATH"/../../bin/migrations

echo >&2 "Run server..."
"$SCRIPTPATH"/../../bin/guideliner
