#!/bin/sh

set -e

MODE=${MODE:-staging}
GOOSE_CMD=${1:-up}

echo "Starting migration for environment: $MODE with command: $GOOSE_CMD"

if [ "$MODE" = "production" ]; then
    DB_HOST="ref-production-postgres"
    DB_PORT="5432"
elif [ "$MODE" = "staging" ]; then
    DB_HOST="ref-staging-postgres" 
    DB_PORT="5432"
else
    echo "Unknown environment: $MODE"
    exit 1
fi

export GOOSE_DRIVER=postgres
export GOOSE_DBSTRING="postgres://postgres:123@${DB_HOST}:${DB_PORT}/refina?sslmode=disable"
export GOOSE_MIGRATION_DIR=/app/migrations

echo "Checking database connection..."
until pg_isready -h $DB_HOST -p $DB_PORT -U postgres; do
    echo "Waiting for database to be ready..."
    sleep 2
done

echo "Database is ready. Running goose $GOOSE_CMD..."
echo "GOOSE_DRIVER: $GOOSE_DRIVER"
echo "GOOSE_DBSTRING: $GOOSE_DBSTRING"
echo "GOOSE_MIGRATION_DIR: $GOOSE_MIGRATION_DIR"

# Run goose command
cd /app
goose "$GOOSE_CMD" "$@"

echo "Goose $GOOSE_CMD completed successfully for $MODE environment"