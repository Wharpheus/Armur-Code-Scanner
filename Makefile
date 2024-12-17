# Run docker compose up
docker-prod:
	@echo "Starting Docker containers in prod mode..."
	docker-compose up --build -d

# Run docker compose up dev
docker-dev:
	@echo "Starting Docker containers in dev mode..."
	docker-compose -f docker-compose.dev.yml up --build -d

# Stop docker containers
docker-down:
	@echo "Stopping Docker containers..."
	docker-compose down

.PHONY: build run test docker-build docker-up docker-down swagger clean dev prod