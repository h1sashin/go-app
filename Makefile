include .env.local

# Configuration
MIGRATIONS_DIR = ./db/migrations
MIGRATE_CMD = migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)"

.PHONY: help migrate rollback force status

help:
	@echo "Usage:"
	@echo "  make run                               # Run development server"
	@echo "  make build                             # Build the application"
	@echo "  make migrate-new name=<migration_name> # Create a new migration"
	@echo "  make migrate-up                        # Apply all pending migrations"
	@echo "  make migrate-down                      # Rollback the last migration"
	@echo "  make migrate-rollback-to version=<v>   # Rollback to a specific version"
	@echo "  make migrate-force version=<v>         # Force a specific migration version"
	@echo "  make migrate-status                    # Show the current migration version"

# Server

# Run development server
run:
	air -c .air.toml

# Build the application
build:
	go build -o bin/app cmd/main.go

# GraphQL

generate-gql:
	go run github.com/99designs/gqlgen generate --config tool/gqlgen/public.gqlgen.yml
	go run github.com/99designs/gqlgen generate --config tool/gqlgen/admin.gqlgen.yml

generate-sql:
	sqlc generate -f tool/sqlc/sqlc.yml

# Migrations

# Create a new migration
migrate-new:
	@test -n "$(name)" || (echo "Please provide a migration name: make migrate new name=<migration_name>"; exit 1)
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(name)

# Apply all pending migrations
migrate-up:
	$(MIGRATE_CMD) up

# Rollback the last migration
migrate-down:
	$(MIGRATE_CMD) down 1

# Rollback to a specific version
migrate-rollback-to:
	@test -n "$(version)" || (echo "Please provide a version: make migrate rollback-to version=<v>"; exit 1)
	$(MIGRATE_CMD) goto $(version)

# Force a specific migration version
migrate-force:
	@test -n "$(version)" || (echo "Please provide a version: make migrate force version=<v>"; exit 1)
	$(MIGRATE_CMD) force $(version)

# Show the current migration version
migrate-status:
	$(MIGRATE_CMD) version
