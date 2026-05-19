include .env
export

export PROJECT_ROOT = $(shell pwd)

env-postgres-up:
	@docker-compose up -d url-postgres

env-postgres-down:
	@docker-compose down url-postgres

env-port-forward:
	@docker-compose up -d port-forward

env-port-close:
	@docker-compose down port-forward

migrate-create:
	$(if $(seq),,$(error seq is not set. Usage: make migrate-create seq=name))
	@docker-compose run --rm url-migrate \
	    create \
	    -ext sql \
	    -dir ./migrations \
	    -seq "$(seq)"

migrate-up:
	@make migrate-action action=up

migrate-down:
	@make migrate-action action=down

migrate-action:
	@docker-compose run --rm url-migrate \
	    -path ./migrations \
	    -database postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable \
	    "$(action)"
