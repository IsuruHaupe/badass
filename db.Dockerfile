FROM mysql:8.0.23

COPY ./database/crete_tables.sql /docker-entrypoint-initdb.d/