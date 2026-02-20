#! /bin/sh

hosts=$(nslookup -type=srv "${SESSION_STORE_HOSTNAME}" | grep -Eo "[^ ]*[0-9]+\.${SESSION_STORE_HOSTNAME}\..*\.svc\..*[^ ]")

max_tries=10
tries=0
pong=0

while [ ${max_tries} -gt ${tries} ]; do
  redis-cli --user "${SESSION_STORE_USERNAME}" --pass "${SESSION_STORE_PASSWORD}" -h "${SESSION_STORE_HOSTNAME}" -p "${SESSION_STORE_PORT}" -t 1 ping
  if [ $? -eq 0 ]; then
    pong=1
    break
  fi

  tries=$((${tries} + 1))
done

addrs=""
if [ ${pong} -eq 1 ]; then
  for host in ${hosts}; do
    addrs="${addrs} ${host}:${SESSION_STORE_PORT}"
  done
  echo $addrs

  echo yes | redis-cli --user "${SESSION_STORE_USERNAME}" --pass "${SESSION_STORE_PASSWORD}" -h "${SESSION_STORE_HOSTNAME}" -p "${SESSION_STORE_PORT}" --cluster create ${addrs}
else
  exit 1
fi
