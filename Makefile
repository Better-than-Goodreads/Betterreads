SHELL := /bin/bash
PWD := $(shell pwd)

default: build

all:


clean-images:
	docker rmi `docker images --filter label=intermediateStageToBeDeleted=true -q`
.PHONY: clean-images

docker-compose-up: docker-image
	docker compose up -d --build
.PHONY: docker-compose-up

docker-compose-down:
	docker compose -f docker-compose-dev.yaml stop -t 1
	docker compose -f docker-compose-dev.yaml down
.PHONY: docker-compose-down

docker-compose-logs:
	docker compose -f docker-compose-dev.yaml logs -f
.PHONY: docker-compose-logs