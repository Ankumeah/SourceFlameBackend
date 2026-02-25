#! /bin/sh

export POSTGRES_USER="${DATABASE_USER}"
export POSTGRES_PASSWORD="${DATABASE_PASSWORD}"
export POSTGRES_DB="${DATABASE_DB}"
export POSTGRES_PORT="${POSTGRES_PORT}"

exec docker-entrypoint.sh postgres
