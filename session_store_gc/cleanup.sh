#! /bin/sh

max_tries=10
tries=0

now=$(date "+%s")

while [ ${max_tries} -gt ${tries} ]; do
  redis-cli \
    --user "${SESSION_STORE_GC_USERNAME}" \
    --pass "${SESSION_STORE_GC_PASSWORD}" \
    -h "${SESSION_STORE_HOSTNAME}" \
    -p "${SESSION_STORE_PORT}" \
    -t 1 \
    zremrangebyscore -inf $now
  if [ $? -eq 0 ]; then
    pong=1
    break
  fi

  tries=$((${tries} + 1))
done

if [ ${pong} -ne 1 ]; then
  exit 1
fi
