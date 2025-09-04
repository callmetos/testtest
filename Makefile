.PHONY: run test lint docker-up docker-down migrate-up migrate-down seed

# Development
run:
	go run ./cmd/main.go

dev: run

# Testing
test:
	go test ./...

# Linting
lint:
	golangci-lint run ./...

# Docker services
docker-up:
	docker compose up -d postgres minio

docker-down:
	docker compose down -v

# Database migrations
migrate-up:
	@echo "Running migrations..."
	@for f in migrations/*.up.sql; do \
		echo "Applying $$f..."; \
		docker exec -i $$(docker compose ps -q postgres) psql -U postgres -d navmate < "$$f"; \
	done

migrate-down:
	@echo "Rolling back migrations..."
	@for f in migrations/*.down.sql; do \
		echo "Rolling back $$f..."; \
		docker exec -i $$(docker compose ps -q postgres) psql -U postgres -d navmate < "$$f"; \
	done

# Seed database
seed:
	@echo "Seeding database..."
	docker exec -i $$(docker compose ps -q postgres) psql -U postgres -d navmate < scripts/seed.sql

# Build
build:
	CGO_ENABLED=0 GOOS=linux go build -o bin/api ./cmd/main.go

# Clean
clean:
	rm -rf bin/