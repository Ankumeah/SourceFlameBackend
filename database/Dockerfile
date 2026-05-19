FROM postgres:alpine

COPY . /opt/postgres/
RUN chmod +x /opt/postgres/run_postgres.sh

ENTRYPOINT [ "/opt/postgres/run_postgres.sh" ]
