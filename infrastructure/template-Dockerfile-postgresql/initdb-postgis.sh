#!/bin/bash

set -e

# Perform all actions as $POSTGRES_USER
export PGUSER="$POSTGRES_USER"

echo "Creating PostGIS template database..."

# Create the 'template_postgis' template db
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --no-password --no-psqlrc <<-EOSQL
	CREATE DATABASE template_postgis IS_TEMPLATE true;
EOSQL

# Load PostGIS into both template_database and $POSTGRES_DB
for DB in template_postgis "$POSTGRES_DB"; do
	echo "Loading PostGIS extensions into $DB"
	psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --no-password --no-psqlrc --dbname="$DB" <<-EOSQL
		-- Create PostGIS extension
		CREATE EXTENSION IF NOT EXISTS postgis;
		CREATE EXTENSION IF NOT EXISTS postgis_topology;

		-- Reconnect to update pg_setting.resetval
		-- See https://github.com/postgis/docker-postgis/issues/288
		\c

		-- Create additional extensions
		CREATE EXTENSION IF NOT EXISTS fuzzystrmatch;
		CREATE EXTENSION IF NOT EXISTS postgis_tiger_geocoder;

		-- Create pg_stat_statements extension
		CREATE EXTENSION IF NOT EXISTS pg_stat_statements;

		-- Optional PostGIS extensions (create if available)
		DO \$\$
		BEGIN
			CREATE EXTENSION IF NOT EXISTS postgis_raster;
		EXCEPTION WHEN others THEN
			RAISE NOTICE 'postgis_raster extension not available';
		END;
		\$\$;

		DO \$\$
		BEGIN
			CREATE EXTENSION IF NOT EXISTS postgis_sfcgal;
		EXCEPTION WHEN others THEN
			RAISE NOTICE 'postgis_sfcgal extension not available';
		END;
		\$\$;
EOSQL
done

echo "PostGIS and pg_stat_statements extensions loaded successfully!"
