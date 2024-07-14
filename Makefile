#!/usr/bin/make
include .env

DB_CONTAINER_IP := $(shell docker inspect -f '{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}' $(POSTGRES_CONTAINER_NAME))

clean: ## Remove images from local registry
	-$(docker_compose_bin) -f docker-compose.yml down
	$(foreach image,$(all_images),$(docker_bin) rmi -f $(image);)

#Actions:
build:
	$(docker_compose_bin) -f docker-compose.yml build

up:
	$(docker_compose_bin) -f docker-compose.yml up -d --build

db-migrations-up:
	GOOSE_DRIVER=postgres GOOSE_MIGRATION_DIR=$(MIGRATION_PATH) GOOSE_DBSTRING="host=$(DB_CONTAINER_IP) user=$(POSTGRES_USER) dbname=$(POSTGRES_DB) sslmode=disable password=$(POSTGRES_PASSWORD)" goose up

db-migrations-down:
	GOOSE_DRIVER=postgres GOOSE_MIGRATION_DIR=$(MIGRATION_PATH) GOOSE_DBSTRING="host=$(DB_CONTAINER_IP) user=$(POSTGRES_USER) dbname=$(POSTGRES_DB) sslmode=disable password=$(POSTGRES_PASSWORD)" goose down

db-migrations-up-by-one:
	GOOSE_DRIVER=postgres GOOSE_MIGRATION_DIR=$(MIGRATION_PATH) GOOSE_DBSTRING="host=$(DB_CONTAINER_IP) user=$(POSTGRES_USER) dbname=$(POSTGRES_DB) sslmode=disable password=$(POSTGRES_PASSWORD)" goose up-by-one
