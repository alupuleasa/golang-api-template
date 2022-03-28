.PHONY: help
help: ## This help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\nTargets:\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)
.DEFAULT_GOAL := help
SHELL:=/bin/bash
BUILD_NAME ?=wallet-service

IMAGENAME=$(BUILD_NAME)
DOCKERCMD=docker exec -t $(IMAGENAME) $(MAKECMD)
DOCKERICMD=docker exec -i -t $(IMAGENAME) $(MAKECMD)

MAKECMD= /usr/bin/make --no-print-directory
RUNCMD=run
GOCMD=go
GOBUILD=$(GOCMD) build -mod=vendor -o wallet-service
GOCLEAN=$(GOCMD) clean -mod=vendor
GOTEST=$(GOCMD) test -mod=vendor -covermode=set -coverprofile=coverage.out ./...
GOCOV=$(GOCMD) tool cover -html=coverage.out -o coverage.html
GOGET=$(GOCMD) get

check: ## Runs different golang checks
	docker run --rm -v $$(pwd):/app -w /app golangci/golangci-lint:v1.23.3 golangci-lint run -E misspell -E dupl -E gocyclo -E golint -E gofmt -E bodyclose -E unparam -E gocritic --modules-download-mode=vendor --timeout 10m ./...
	
shell: 
	@docker-compose exec wallet-service bash

test: ## Runs go tests
	@if [ -f /.dockerenv ] ; then \
		$(GOTEST); \
	else \
		echo "running tests in docker env ..." ; \
		$(DOCKERCMD) test;\
	fi;
	@exit $(.SHELLSTATUS)

build: ## builds the docker image
	@if [ -f /.dockerenv ] ; then \
		$(GOBUILD); \
	else \
	  echo "building app in docker env ..." ; \
	  $(DOCKERCMD) build; \
  	fi;

run: 
	@if [ -f /.dockerenv ] ; then \
		pkill $(BUILD_NAME);\
		./$(BUILD_NAME) $(RUNCMD) $@ ;\
	else \
		echo "running app in docker env ..." ; \
		$(DOCKERICMD) run; \
	fi;

brun: ## builds and runs the new binary
	@$(MAKECMD) build
	@$(MAKECMD) run

up: ## Brings the docker images up
	@[ -f /.dockerenv ] || true
	@docker-compose pull
	@docker network create wallet-network --subnet=192.168.100.0/24 || true
	@docker-compose up -d

reload: ## Brings the docker images up
	@[ -f /.dockerenv ] || true
	@docker-compose down --volumes
	@docker-compose up --force-recreate --build -d -V

reload_debug: ## Brings the docker images up
	@[ -f /.dockerenv ] || true
	@docker-compose down --volumes
	@docker-compose up --force-recreate --build -V

stop: ## Stops the docker images
	@[ -f /.dockerenv ] || true
	@docker-compose stop

pause:  ## Pauses the docker images
	@[ -f /.dockerenv ] || true
	@docker-compose pause

restart: ## Restarts the docker images
	@[ -f /.dockerenv ] || true
	@docker-compose restart

start: ## Starts the docker images
	@[ -f /.dockerenv ] || true
	@docker-compose start

down:  ## Destroys the docker images
	@[ -f /.dockerenv ] || true
	@docker-compose down --remove-orphans -v

docker_clean:  ## Destroys the docker images
	@[ -f /.dockerenv ] || true
	@docker-compose down --volumes --remove-orphans
	@docker-compose rm --force -v

ps:  ## Shows the docker image stats
	@[ -f /.dockerenv ] || true
	@docker ps -a
