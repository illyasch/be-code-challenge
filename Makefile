.PHONY: \
	image \
	up \
	docker-compose \
	test \
	\

.DEFAULT_GOAL:=help

SHELL := /bin/bash

# ==============================================================================
# Building containers

# $(shell git rev-parse --short HEAD)
VERSION := dev
DOCKER_COMPOSE_FILE := infra/docker-compose.yml

image:		## Build a challenge service image in docker
	@:$(call check_defined, VERSION, version)
	docker build \
		-f infra/challenge.Dockerfile \
		-t challenge:$(VERSION) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

up:             ## Build and start challenge service and dependencies
	docker-compose -f $(DOCKER_COMPOSE_FILE) up --build -d --force-recreate --remove-orphans challenge

clean:          ## Stops running services and removes containers, volumes and images
	docker-compose -f $(DOCKER_COMPOSE_FILE) kill
	docker-compose -f $(DOCKER_COMPOSE_FILE) rm -sfv

# ==============================================================================
# Running tests within the local computer

test: 		## Run tests inside a container
	docker-compose -f $(DOCKER_COMPOSE_FILE) build test
	docker-compose -f $(DOCKER_COMPOSE_FILE) run --rm test

# ==============================================================================
# Help

help:		## Show this help message
	@echo
	@echo '  Usage:'
	@echo '    make <target>'
	@echo
	@echo '  Targets:'
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'
	@echo
