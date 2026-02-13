#! /bin/sh

envsubst '$SESSION_STORE_USERNAME #SESSION_STORE_PASSWORD' < /opt/redis/redis.conf.template > /opt/redis/redis.conf

exec "/usr/local/bin/docker-entrypoint.sh" /opt/redis/redis.conf --port "${SESSION_STORE_PORT}"
