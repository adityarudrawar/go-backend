1. Create the custom postgres image by running the db.Dockerfile
```
Build the image Command - docker image build -f db.Dockerfile . -t postgresdb:latest
```

If you want to run the image, the command is
```
docker run --name metalman66/postgresdb -e POSTGRES_PASSWORD=postgres -p 5432:5432 -h postgres -d postgres
```

2. Optional: Upload to docker hub

3. Run the Docker compose file
```
docker-compose -f docker-compose.yml up -d
```

4. The application should be running on port 8085

#### Packages that I am using
1. ORM: GROM: go get gorm.io/gorm
2. POSTGRES DRIVER:  go get gorm.io/driver/postgres
https://gorm.io/docs/connecting_to_the_database.html
3. 

#### Useful links
1. https://cadu.dev/creating-a-docker-image-with-database-preloaded/
2. https://stackoverflow.com/questions/53735948/cannot-get-connection-between-containers
