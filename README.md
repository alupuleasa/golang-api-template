# Part 1: API Service
(part 2 at the bottom)
&nbsp; 

# Running the service in docker env
```
docker compose up -d
make shell
make build
make run # or ./api run
```

&nbsp; 
# Tests
## Running tests in docker env
```
docker compose up -d
make shell
make test
```

&nbsp; 


# Configure service for other envs see --help option 
```
./api run --help
```

You can use args (1st column name) or environment (2nd column name) to update these running parameters. [Docs](https://github.com/synthesio/zconfig)

&nbsp; 

# Fixtures
## Injected postgres fixtures to the postgres docker container in the ./db folder which are loaded on the `docker compose up` first run (or volume clear)
&nbsp; 
# Misc
### Reload containers and database 
```
docker compose down -v && docker compose up -V -d
```

&nbsp; 

