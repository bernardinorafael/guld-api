include .env
default: run
# Sets variable for common migration Docker command
MIGRATE_CMD = docker run -it --rm --network host --volume $(PWD)/internal/infra/database/pg:/db migrate/migrate

# Access the PostgreSQL container
pg-container:
	@docker exec -it $(DB_NAME) psql -U $(DB_USER) -d $(DB_NAME)
.PHONY: pg-container

# Execute the Go server
run:
	@go run cmd/server/main.go
.PHONY: run

# Add a new migration
migration:
	@if [ -z "$(name)" ]; then echo "Migration name is required"; exit 1; fi
	@$(MIGRATE_CMD) create -ext sql -dir /db/migrations $(name)
.PHONY: migration

# Apply all pending migrations
migrate-up:
	@$(MIGRATE_CMD) -path=/db/migrations -database "$(DB_URL)" up
.PHONY: up

# Revert all applied migrations
migrate-down:
	@$(MIGRATE_CMD) -path=/db/migrations -database "$(DB_URL)" down
.PHONY: down

# Apply the last pending migration
migrate-next:
	@$(MIGRATE_CMD) -path=/db/migrations -database "$(DB_URL)" up 1
.PHONY: migrate-next

# Revert the last applied migration
migrate-previous:
	@$(MIGRATE_CMD) -path=/db/migrations -database "$(DB_URL)" down 1
.PHONY: migrate-previous
