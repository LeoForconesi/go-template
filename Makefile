SHELL := /bin/bash

APP_NAME := go-template
COMPOSE := docker compose

.PHONY: help
help:
	@echo "Targets:"
	@echo "  make up           - Levanta todo (infra + app)"
	@echo "  make down         - Baja todo"
	@echo "  make logs         - Sigue logs de app"
	@echo "  make ps           - Lista servicios"
	@echo "  make build        - Build de la imagen app"
	@echo "  make reset        - down -v + up limpio"
	@echo "  make migrate-up   - Aplica migraciones (contenedor migrate)"
	@echo "  make migrate-down - Revierte migraciones (1 step)"
	@echo "  make kafka-ui 	   - Levanta Redpanda Console para ver topics de Kafka en http://localhost:8081"
	@echo "  make test         - go test ./..."
	@echo "  make tidy         - go mod tidy"

.PHONY: up
up:
	$(COMPOSE) up -d --build

.PHONY: down
down:
	$(COMPOSE) down

.PHONY: logs
logs:
	$(COMPOSE) logs -f app

.PHONY: ps
ps:
	$(COMPOSE) ps

.PHONY: build
build:
	docker build -t $(APP_NAME):local .

.PHONY: reset
reset:
	$(COMPOSE) down -v
	$(COMPOSE) up -d --build

# Ejecuta el contenedor migrate manualmente (por si necesitás re-aplicar)
.PHONY: migrate-up
migrate-up:
	$(COMPOSE) run --rm migrator

.PHONY: migrate-down
migrate-down:
	$(COMPOSE) run --rm migrator sh -c "migrate -path /migrations -database 'postgres://app:app@postgres:5432/appdb?sslmode=disable' down 1"

kafka-ui:
	@echo "🚀 Levantando Redpanda Console en http://localhost:8081..."
	docker compose up -d redpanda-console
	@sleep 3
	@open http://localhost:8081 || xdg-open http://localhost:8081 || echo "Abierto en http://localhost:8081"

.PHONY: test
test:
	go test ./... -v

.PHONY: tidy
tidy:
	go mod tidy

GOBIN ?= $$(go env GOPATH)/bin

.PHONY: install-go-test-coverage
install-go-test-coverage:
	go install github.com/vladopajic/go-test-coverage/v2@latest

.PHONY: check-coverage
check-coverage: install-go-test-coverage
	go test ./... -coverprofile=./cover.out -covermode=atomic -coverpkg=./...
	${GOBIN}/go-test-coverage --config=testcoverage.yml