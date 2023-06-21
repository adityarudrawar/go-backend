1. Create the custom postgres image by running the db.Dockerfile
```
Build the image Command - docker image build -f db.Dockerfile . -t nimblepostgresdb:latest
```

If you want to run the image

```
Run Command - docker container run -d --rm -p 5432:5432 -e POSTGRES_PASSWORD=postgres -e POSTGRES_USER=postgres --name nimblepostgresdb nimblepostgresdb:latest
```

2. Optional: Upload to docker hub

3. Run the Docker compose file
```
docker-compose -f docker-compose.yml up -d
```
4. The application should be running on port 8085