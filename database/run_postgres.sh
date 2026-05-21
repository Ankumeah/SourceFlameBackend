#! /bin/sh

export POSTGRES_USER=$(echo "$DATABASE_URL" | awk -F/ '{print $3}' | awk -F: '{print $1}')
export POSTGRES_PASSWORD=$(echo "$DATABASE_URL" | awk -F/ '{print $3}' | awk -F: '{print $2}' | awk -F@ '{print $1}')
export POSTGRES_PORT=$(echo "$DATABASE_URL" | awk -F/ '{print $3}' | awk -F: '{print $3}')
export POSTGRES_DB=$(echo "$DATABASE_URL" | awk -F/ '{print $4}' | awk -F\? '{print $1}')

exec docker-entrypoint.sh postgres
