#!/bin/bash

NAME=$1
if [ -z "$NAME" ]; then
  echo "Usage: ./create_migration.sh <migration_name>"
  exit 1
fi

TIMESTAMP=$(date +"%Y%m%d%H%M%S")
mkdir -p migrations
touch "migrations/${TIMESTAMP}_${NAME}.up.sql"
touch "migrations/${TIMESTAMP}_${NAME}.down.sql"

echo "Created:"
echo " - migrations/${TIMESTAMP}_${NAME}.up.sql"
echo " - migrations/${TIMESTAMP}_${NAME}.down.sql"
