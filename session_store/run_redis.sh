#! /bin/sh

for i in cluster standalone; do
  sed -i "s/SESSIONS_USERNAME/$SESSION_STORE_SESSIONS_USERNAME/g" /opt/redis/redis_$i.conf
  sed -i "s/SESSIONS_PASSWORD/$SESSION_STORE_SESSIONS_PASSWORD/g" /opt/redis/redis_$i.conf

  sed -i "s/GC_USERNAME/$SESSION_STORE_GC_USERNAME/g" /opt/redis/redis_$i.conf
  sed -i "s/GC_PASSWORD/$SESSION_STORE_GC_PASSWORD/g" /opt/redis/redis_$i.conf
done

if test "$SESSION_STORE_TYPE" = "redis_cluster"; then
  exec "/usr/local/bin/docker-entrypoint.sh" /opt/redis/redis_cluster.conf --port "${SESSION_STORE_PORT}"
elif test "$SESSION_STORE_TYPE" = "redis_standalone"; then
  exec "/usr/local/bin/docker-entrypoint.sh" /opt/redis/redis_standalone.conf --port "${SESSION_STORE_PORT}"
else
  echo "Unsupported session store type $SESSION_STORE_TYPE"
  exit 1
fi
