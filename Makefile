# Run docker compose up dev
docker-up:
	@echo "Starting Docker containers in dev mode..."
	docker-compose up --build -d

# Stop docker containers
docker-down:
	@echo "Stopping Docker containers..."
	docker-compose down

.PHONY: build run test docker-build docker-up docker-down swagger clean dev prod