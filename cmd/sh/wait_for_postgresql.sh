#!/bin/sh

echo >&2 "Wait for postgresql..."
until PGPASSWORD=$GUIDELINER_DB_PASSWORD psql \
    -h "$GUIDELINER_DB_HOST" \
    -U "$GUIDELINER_DB_LOGIN" \
    -d "$GUIDELINER_DB_NAME" \
    -p "$GUIDELINER_DB_PORT" \
    -c '\q'; do
  echo >&2 "Postgres is unavailable - sleeping"
  sleep 1
done
echo >&2 "Postgres is up - executing command"
