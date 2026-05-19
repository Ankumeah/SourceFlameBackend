#! /bin/sh

export POSTGRES_USER=$(echo "$DATABASE_URL" | awk -F/ '{print $3}' | awk -F: '{print $1}')
export POSTGRES_PASSWORD=$(echo "$DATABASE_URL" | awk -F/ '{print $3}' | awk -F: '{print $2}' | awk -F@ '{print $1}')
export POSTGRES_PORT=$(echo "$DATABASE_URL" | awk -F/ '{print $3}' | awk -F: '{print $3}')
export POSTGRES_DB=$(echo "$DATABASE_URL" | awk -F/ '{print $4}' | awk -F\? '{print $1}')
echo POSTGRES_DB $POSTGRES_DB POSTGRES_PORT $POSTGRES_PORT POSTGRES_PASSWORD $POSTGRES_PASSWORD POSTGRES_USER $POSTGRES_USER

exec docker-entrypoint.sh postgres
