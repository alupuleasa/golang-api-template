# Part 1: Wallet Service
Service whose sole responsibility is to store the
funds and provide the functionality for manipulating the balance.

&nbsp; 

## [API documentation](https://documenter.getpostman.com/view/2833399/TzeWF7TQ)

&nbsp; 

# Running the service in docker env
```
docker compose up -d
make shell
make build
make run # or ./wallet run
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
./wallet run --help
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


# Part 2: [PayPal workflow](https://www.websequencediagrams.com/files/render?link=bqEqgYcgtmXAbBDod1BQGGiye84xcKO5uRR0kIoBMEt1i7FxvpsRrns9dT84lqaU)
