FROM redis:alpine

COPY . /opt/redis/
RUN chmod +x /opt/redis/run_redis.sh

EXPOSE 6379

ENTRYPOINT [ "/opt/redis/run_redis.sh" ]
