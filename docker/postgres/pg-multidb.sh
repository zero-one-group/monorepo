#!/usr/bin/env bash
set -euo pipefail

###
# Example usage on the compose.yaml file:
#
# env:
#   PG_DATABASES: mydb1,mydb2
###

# Parse and process multiple databases and schemas
if [ -n "${PG_DATABASES:-}" ]; then
    echo "Multiple database creation requested: $PG_DATABASES"
    IFS=',' read -ra databases <<< "$PG_DATABASES"

    # Create extra databases with owners
    for db in "${databases[@]}"; do
        echo "  Creating database '$db' with owner '$POSTGRES_USER'"
        psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
            CREATE DATABASE "$db" WITH OWNER = "$POSTGRES_USER"
            ENCODING = 'UTF8' LC_COLLATE = 'en_US.UTF-8'
            LC_CTYPE = 'en_US.UTF-8'
            TEMPLATE = template0;
EOSQL
    done
    echo "Multiple databases and schemas created successfully."
else
    echo "No databases requested for creation."
fi
