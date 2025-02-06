#!/bin/bash
###
# Example usage on the docker-compose.yaml file:
#
# env:
#   POSTGRES_EXTRA_DB: mydb1,mydb2
#   POSTGRES_SCHEMAS: mydb1:newschema,mydb2:myschema
###

set -e
set -u

# Function to create a database with owner
function create_database_with_owner() {
    local database_name=$1
    local owner_name=$2
    echo "  Creating database '$database_name' with owner '$owner_name'"
    psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" <<-EOSQL
        CREATE DATABASE "$database_name"
        WITH OWNER = "$owner_name"
        ENCODING = 'UTF8'
        LC_COLLATE = 'en_US.UTF-8'
        LC_CTYPE = 'en_US.UTF-8'
        TEMPLATE = template0;
EOSQL
}

# Function to create schema and set ownership and permissions
function create_schema_with_permissions() {
    local database_name=$1
    local schema_name=$2
    local schema_owner=$3
    echo "  Creating schema '$schema_name' in database '$database_name' with owner '$schema_owner'"
    psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname="$database_name" <<-EOSQL
        DO \$\$
        DECLARE
            sch_name TEXT := '$schema_name';
            sch_owner TEXT := '$schema_owner';
        BEGIN
            -- Check if the schema exists
            IF NOT EXISTS (
                SELECT 1
                FROM information_schema.schemata
                WHERE schema_name = sch_name
            ) THEN
                -- Create the schema
                EXECUTE format('CREATE SCHEMA %I AUTHORIZATION %I', sch_name, sch_owner);
            END IF;

            -- Grant permissions to the owner
            EXECUTE format('GRANT ALL PRIVILEGES ON SCHEMA %I TO %I', sch_name, sch_owner);

            -- Allow the owner full access to future objects in the schema
            EXECUTE format('ALTER DEFAULT PRIVILEGES IN SCHEMA %I GRANT ALL PRIVILEGES ON TABLES TO %I', sch_name, sch_owner);
            EXECUTE format('ALTER DEFAULT PRIVILEGES IN SCHEMA %I GRANT ALL PRIVILEGES ON SEQUENCES TO %I', sch_name, sch_owner);
            EXECUTE format('ALTER DEFAULT PRIVILEGES IN SCHEMA %I GRANT ALL PRIVILEGES ON FUNCTIONS TO %I', sch_name, sch_owner);
        END
        \$\$;
EOSQL
}

# Parse and process multiple databases and schemas
if [ -n "${POSTGRES_EXTRA_DB:-}" ]; then
    echo "Multiple database creation requested: $POSTGRES_EXTRA_DB"
    IFS=',' read -ra databases <<< "$POSTGRES_EXTRA_DB"
    declare -A schema_map

    # Parse additional schemas if provided
    if [ -n "${POSTGRES_SCHEMAS:-}" ]; then
        echo "Processing additional schemas: $POSTGRES_SCHEMAS"
        IFS=',' read -ra schema_entries <<< "$POSTGRES_SCHEMAS"
        for entry in "${schema_entries[@]}"; do
            db=$(echo "$entry" | cut -d':' -f1)
            schema=$(echo "$entry" | cut -d':' -f2)
            schema_map["$db"]="$schema"
        done
    fi

    # Create the databases
    for db in "${databases[@]}"; do
        create_database_with_owner "$db" "$POSTGRES_USER"
    done

    # Create schemas after databases
    for db in "${databases[@]}"; do
        if [ -n "${schema_map[$db]:-}" ]; then
            create_schema_with_permissions "$db" "${schema_map[$db]}" "$POSTGRES_USER"
        fi
    done
    echo "Multiple databases and schemas created successfully."
else
    echo "No databases requested for creation."
fi
