#!/bin/bash
# This script is run when the PostgreSQL container is first created.
# It creates an additional database for testing purposes.
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE DATABASE "postgres-test";
EOSQL