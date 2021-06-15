# wallet




## Running tests
```
docker compose up -d
make shell
make test
```

# Running the service in docker
```
docker compose up -d
make shell
make build
./wallet
```

## Added fixtures to the docker compose container in the ./db folder
### Reload containers and database 
```
docker compose down -v && docker compose up -V -d
```