#! /bin/sh

envsubst '$BACKEND_HOSTNAME $BACKEND_PORT' < /opt/nginx/nginx.conf.template > /opt/nginx/nginx.conf

exec nginx -g "daemon off;" -c /opt/nginx/nginx.conf
