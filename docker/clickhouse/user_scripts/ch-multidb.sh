#!/usr/bin/env bash
set -euo pipefail

###
# Example usage on the compose.yaml file:
#
# env:
#   CLICKHOUSE_DATABASES: mydb1,mydb2
###

# Parse and process multiple databases and schemas
if [ -n "${CLICKHOUSE_DATABASES:-}" ]; then
    echo "Multiple database creation requested: $CLICKHOUSE_DATABASES"
    IFS=',' read -ra databases <<< "$CLICKHOUSE_DATABASES"

    # Create extra databases
    for db in "${databases[@]}"; do
        echo "  Creating database '$db'"
        clickhouse client -q "CREATE DATABASE \"$db\""
    done

    echo "Multiple databases created successfully."
else
    echo "No databases requested for creation."
fi
