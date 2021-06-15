.PHONY: help
help: ## This help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\nTargets:\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)
.DEFAULT_GOAL := help
SHELL:=/bin/bash
BUILD_NAME ?=catchon-serve
GOCMD=go
GOBUILD=$(GOCMD) build -mod=vendor -o wallet
GOCLEAN=$(GOCMD) clean -mod=vendor
GOTEST=$(GOCMD) test -mod=vendor -covermode=set -coverprofile=coverage.out ./...
GOCOV=$(GOCMD) tool cover -html=coverage.out -o coverage.html
GOGET=$(GOCMD) get

check: ## Runs different golang checks
	docker run --rm -v $$(pwd):/app -w /app golangci/golangci-lint:v1.23.3 golangci-lint run -E misspell -E dupl -E gocyclo -E golint -E gofmt -E bodyclose -E unparam -E gocritic --modules-download-mode=vendor --timeout 10m ./...
	
shell: 
	@docker compose exec wallet-service bash

test: ## Runs go tests
	@[ -f /.dockerenv ] || (echo "not in a docker environment. run > make shell"  && false)
	$(GOTEST); 

build: ## builds the docker image
	@[ -f /.dockerenv ] || (echo "not in a docker environment. run > make shell"  && false)
	$(GOBUILD); 

run: 
	@make build
	./wallet run


up: ## Brings the docker images up
	@[ -f /.dockerenv ] || true
	@docker compose up -d

reload: ## Brings the docker images up
	@[ -f /.dockerenv ] || true
	@docker compose down --volumes
	@docker compose up --force-recreate --build -d -V

reload_debug: ## Brings the docker images up
	@[ -f /.dockerenv ] || true
	@docker compose down --volumes
	@docker compose up --force-recreate --build -V

stop: ## Stops the docker images
	@[ -f /.dockerenv ] || true
	@docker compose stop

pause:  ## Pauses the docker images
	@[ -f /.dockerenv ] || true
	@docker compose pause

restart: ## Restarts the docker images
	@[ -f /.dockerenv ] || true
	@docker compose restart

start: ## Starts the docker images
	@[ -f /.dockerenv ] || true
	@docker compose start

down:  ## Destroys the docker images
	@[ -f /.dockerenv ] || true
	@docker compose down --remove-orphans -v

docker_clean:  ## Destroys the docker images
	@[ -f /.dockerenv ] || true
	@docker compose down --volumes --remove-orphans
	@docker compose rm --force -v

ps:  ## Shows the docker image stats
	@[ -f /.dockerenv ] || true
	@docker ps -a
