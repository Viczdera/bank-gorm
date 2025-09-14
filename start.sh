#!/bin/sh

set -e

echo "sh command: run db migration"
/app/migrate -path /app/migrations  -database "$DB_SOURCE" -verbose up

echo "sh command: start app"
exec "$@"

