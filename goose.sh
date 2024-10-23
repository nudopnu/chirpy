#!/bin/sh

set -a
GOOSE_DRIVER="postgres"
GOOSE_MIGRATION_DIR="./sql/schema"
. ./.env
GOOSE_DBSTRING="postgres://${DB_USERNAME}:${DB_PASSWORD}@localhost:5432/chirpy"
set +a

goose "$@"