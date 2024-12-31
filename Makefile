include .envrc

MIGRATIONS_PATH = ./cmd/migrate/migrations

.PHONY: migrate-create
migration:
	@migrate create -seq -ext sql -dir $(MIGRATIONS_PATH) $(filter-out $@, $(MAKECMDGOALS))
.PHONY: migrate-up
migrate-up:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) up
.PHONY: migrate-down
migrate-down:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) down
.PHONY: migrate-force
migrate-force:
	@if [ -z "$(filter-out $@, $(MAKECMDGOALS))" ]; then echo "Error: VERSION is required for migrate-force. Try to run: make migrate-force VERSION"; \
	exit 1; \
	fi
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) force $(filter-out $@, $(MAKECMDGOALS))
%:
	@:
.PHONY: migrate-version
migrate-version:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) version
.PHONY: migrate-goto
migrate-goto:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) goto $(filter-out $@, $(MAKECMDGOALS))