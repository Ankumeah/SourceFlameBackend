#! /bin/sh

sed -i "s/SESSION_STORE_USERNAME/${SESSION_STORE_USERNAME}/g" /opt/redis/redis.conf
sed -i "s/SESSION_STORE_PASSWORD/${SESSION_STORE_PASSWORD}/g" /opt/redis/redis.conf

exec "/usr/local/bin/docker-entrypoint.sh" /opt/redis/redis.conf --port "${SESSION_STORE_PORT}"
