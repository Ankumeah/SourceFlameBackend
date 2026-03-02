#! /bin/sh

sed -i "s/SESSIONS_USERNAME/$SESSION_STORE_SESSIONS_USERNAME/g" /opt/redis/redis.conf
sed -i "s/SESSIONS_PASSWORD/$SESSION_STORE_SESSIONS_PASSWORD/g" /opt/redis/redis.conf

sed -i "s/GC_USERNAME/$SESSION_STORE_GC_USERNAME/g" /opt/redis/redis.conf
sed -i "s/GC_PASSWORD/$SESSION_STORE_GC_PASSWORD/g" /opt/redis/redis.conf

exec "/usr/local/bin/docker-entrypoint.sh" /opt/redis/redis.conf --port "${SESSION_STORE_PORT}"
