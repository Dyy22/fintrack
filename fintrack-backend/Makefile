APP_NAME=fintrack-backend
DB_URL=postgres://fintrack:fintrack@localhost:5432/fintrack?sslmode=disable
COMPOSE=docker compose -f docker-compose.dev.yml

.PHONY: help run tidy test test-integration smoke-test dev db-up db-migrate db-reset db-down migrate-up migrate-down docker-up docker-down logs

help:
	@echo "Available targets:"
	@echo "  make run              Run backend locally with Go"
	@echo "  make tidy             Run go mod tidy"
	@echo "  make test             Run unit tests"
	@echo "  make test-integration Run PostgreSQL integration tests"
	@echo "  make dev              Start Postgres/backend and run migrations"
	@echo "  make smoke-test       Run API smoke test against localhost:8080"
	@echo "  make db-up            Start Postgres only"
	@echo "  make db-migrate       Apply initial migration using psql in Docker"
	@echo "  make db-reset         Reset database using migration down/up"
	@echo "  make docker-up        Start all Docker Compose services"
	@echo "  make docker-down      Stop Docker Compose services"
	@echo "  make logs             Follow backend logs"

run:
	go run ./cmd/api

tidy:
	go mod tidy

test:
	go test ./...

test-integration:
	FINTRACK_TEST_DB_URL="$${FINTRACK_TEST_DB_URL:-postgres://fintrack:fintrack@localhost:5432/fintrack_test?sslmode=disable}" go test -tags=integration ./internal/repository/postgres

smoke-test:
	./scripts/smoke-test.sh

dev: db-up db-migrate docker-up

db-up:
	$(COMPOSE) up -d postgres
	$(COMPOSE) exec -T postgres sh -c 'until pg_isready -U fintrack -d fintrack; do sleep 1; done'

db-migrate:
	./scripts/db-migrate.sh

db-reset: db-up
	-$(COMPOSE) exec -T postgres psql -U fintrack -d fintrack -v ON_ERROR_STOP=1 < migrations/000001_init.down.sql
	$(COMPOSE) exec -T postgres psql -U fintrack -d fintrack -v ON_ERROR_STOP=1 < migrations/000001_init.up.sql

db-down:
	$(COMPOSE) down -v

migrate-up:
	migrate -path migrations -database "$(DB_URL)" up

migrate-down:
	migrate -path migrations -database "$(DB_URL)" down 1

docker-up:
	$(COMPOSE) up -d --build

docker-down:
	$(COMPOSE) down

logs:
	$(COMPOSE) logs -f backend
