# -- Build Command - docker image build -f db.Dockerfile . -t metalman66/postgresdb:latest
# -- Run Command - docker container run -d --rm -p 5432:5432 -e POSTGRES_PASSWORD=postgres -e POSTGRES_USER=postgres --name postgresdb metalman66/postgresdb:latest
FROM postgres
ENV POSTGRES_USER=postgres
ENV POSTGRES_PASSWORD=postgres
COPY init.sql /docker-entrypoint-initdb.d/