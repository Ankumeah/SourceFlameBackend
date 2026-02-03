#! /bin/sh

sed -i "s/USERNAME/${USERNAME}/g" /opt/redis/redis.conf
sed -i "s/PASSWORD/${PASSWORD}/g" /opt/redis/redis.conf

exec redis-server /opt/redis/redis.conf
