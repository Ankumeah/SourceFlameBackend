#! /bin/sh

sed -i "s/SESSION_USER_USERNAME/${SESSION_USER_USERNAME}/g" /opt/redis/redis.conf
sed -i "s/SESSION_USER_PASSWORD/${SESSION_USER_PASSWORD}/g" /opt/redis/redis.conf

sed -i "s/REFRESH_USER_USERNAME/${REFRESH_USER_USERNAME}/g" /opt/redis/redis.conf
sed -i "s/REFRESH_USER_PASSWORD/${REFRESH_USER_PASSWORD}/g" /opt/redis/redis.conf

exec redis-server /opt/redis/redis.conf
