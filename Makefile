.PHONY: up down build clean logs test db-up db-down local-run local-build docker-build docker-run docker-stop

# Binary name for local development
BINARY_NAME=payments

# Docker service names
DB_SERVICE_NAME=postgres-db
APP_SERVICE_NAME=payments

# Start all services with Docker Compose for development
up:
	docker-compose up --build -d

# Stop and remove containers
down:
	docker-compose down

# Rebuild services using Docker
docker-build:
	docker-compose build

# Run the application in Docker
docker-run:
	docker-compose up payments

# Stop the Docker application
docker-stop:
	docker-compose stop payments

# Follow logs from Docker containers
logs:
	docker-compose logs -f

# Start only the PostgreSQL container with environment variables for Docker
db-up:
	docker run --name $(DB_SERVICE_NAME) -d \
		-e POSTGRES_USER=postgres \
		-e POSTGRES_PASSWORD=postgres \
		-e POSTGRES_DB=payments \
		-p 5432:5432 \
		-v ./init-scripts:/docker-entrypoint-initdb.d \
		postgres:12

# Stop and remove the PostgreSQL Docker container
db-down:
	docker stop $(DB_SERVICE_NAME)
	docker rm $(DB_SERVICE_NAME)

# Local development commands
local-build:
	cd payments && \
	echo "Building the local development binary..." && \
	go build -o ./bin/$(BINARY_NAME)

local-run: local-build
	cd payments && \
	echo "Running the local development server..." && \
	./bin/$(BINARY_NAME)

test:
	cd payments && \
	echo "Running tests locally..." && \
	go test ./...

# Clean up Docker images and volumes
clean:
	docker system prune -af --volumes