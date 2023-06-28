This is a template that I developed for REST API backend with Go programming language and postgres database.
This has support for custom postgres images that you can create and preload data in to that image, plus already written kubernetes files.

You can directly fork this and update according to your needs.

I am constantly learning and will be improving the code base even more.

## Technologies and libraries
1. [Go](https://go.dev/)
2. [Postgres](https://www.postgresql.org/)
3. [Gofiber](https://gofiber.io/)
4. [ORM](https://gorm.io/)
5. [Docker](https://www.docker.com/)

## Building and running the application

### Create Custom postgres docker image
Modiy the database/init.sql according to your needs to create your databases, tables, relations, and load data.
During the building phase the Docker image will be builded by running the script inside the container to create the custom postgres image.

Command to build the postgres docker image
```
docker image build -f db.Dockerfile -t metalman66/postgresdb:latest .
```

Push the postgres docker image to docker hub
```
docker push metalman66/postgresdb
```

### Docker Compose

You can run the application throught docker by using the docker-compose file.
Command to run docker-compose.

```
docker-compose up --build
```

### Kubernetes
You can run the application through deployment on a kubernetes cluster. 
Yaml files are present for the gobackend and the postgres db.


Run the following command to start the Kubernetes deployment and services.
```
kubectl apply -f postgresdb.yaml
kubectl apply -f gobackend.yaml
```
