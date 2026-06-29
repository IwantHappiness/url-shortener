include .env
export

export PROJECT_ROOT = $(shell pwd)

run-app:
	@export LOGGER_FOLDER=${PROJECT_ROOT}/backend/out/logs && \
	export POSTGRES_HOST=localhost && \
	go mod tidy && \
	go run ${PROJECT_ROOT}/backend/cmd/api/main.go

env-up:
	@docker-compose up -d url-postgres

env-down:
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
	    -dir /migrations \
	    -seq "$(seq)"

migrate-up:
	@make migrate-action action=up

migrate-down:
	@make migrate-action action=down

migrate-action:
	@docker-compose run --rm url-migrate \
	    -path /migrations \
	    -database postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@url-postgres:$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable \
	    "$(action)"

swagger-gen:
	@docker-compose run --rm swagger

up:
	@docker-compose up -d backend frontend

down:
	@docker-compose down

clean-up:
	@read -p "Очистить все volume файлы окружения? Опасность утери данных. [Y/n] :" ans; \
	if [ "$$ans" = "y" ]; then \
		docker-compose down url-postgres port-forward && \
		rm -rf ${PROJECT_ROOT}/backend/out/pgdata && \
		echo "Файлы окружения очищены."; \
	else \
		echo "Очистка окружения отменена."; \
	fi

clean-up-logs:
	@read -p "Очистить все logs файлы? Опасность утери логов. [Y/n] :" ans; \
	if [ "$$ans" = "y" ]; then \
		rm -rf ${PROJECT_ROOT}/backend/out/logs && \
		echo "Файлы логов очищены."; \
	else \
		echo "Очистка логов отменена."; \
	fi
