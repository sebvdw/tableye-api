#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE EXTENSION IF NOT EXISTS pgtle;
    CREATE EXTENSION IF NOT EXISTS pg_tde;
EOSQL

# If the above fails, log the error
if [ $? -ne 0 ]; then
    echo "Error: Failed to create extensions. Check if pgtle and pg_tde are properly installed."
    exit 1
fi

echo "Initialization completed successfully."
