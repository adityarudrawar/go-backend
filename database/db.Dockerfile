# -- Build Command - docker image build -f db.Dockerfile . -t postgresdb:latest
# -- Run Command - docker container run -d --rm -p 5432:5432 -e POSTGRES_PASSWORD=postgres -e POSTGRES_USER=postgres --name postgresdb postgresdb:latest

FROM postgres:11-alpine as dumper

COPY init.sql /docker-entrypoint-initdb.d/

RUN ["sed", "-i", "s/exec \"$@\"/echo \"skipping...\"/", "/usr/local/bin/docker-entrypoint.sh"]

ENV POSTGRES_USER=postgres
ENV POSTGRES_PASSWORD=postgres
ENV PGDATA=/data

RUN ["/usr/local/bin/docker-entrypoint.sh", "postgres"]

# final build stage
FROM postgres:11-alpine

COPY --from=dumper /data $PGDATA