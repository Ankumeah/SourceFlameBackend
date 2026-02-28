FROM postgres:alpine

ARG version=alpha-01
COPY ${version}.sql /docker-entrypoint-initdb.d/

COPY . /opt/postgres/
RUN chmod +x /opt/postgres/run_postgres.sh

ENTRYPOINT [ "/opt/postgres/run_postgres.sh" ]
