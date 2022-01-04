FROM mysql:8.0.23

COPY ./database/*.sql /docker-entrypoint-initdb.d/