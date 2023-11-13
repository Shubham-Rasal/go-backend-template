#!/bin/sh

set -e

# Run db migrations
echo "Running db migrations"
source /app/app.env
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

# Run the main application
echo "Running the main application"
exec "$@"
