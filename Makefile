.PHONY: dev build seed test docker-up docker-down docker-build frontend backend

dev:
	cd backend && go run cmd/main.go

frontend-dev:
	npm run dev

backend:
	cd backend && go run cmd/main.go

build:
	cd backend && go build -o bin/pos-backend ./cmd/main.go
	npm run build

seed:
	curl -X POST http://localhost:8080/api/seed

test:
	cd backend && go test ./...

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

docker-build:
	docker-compose up -d --build

docker-rebuild:
	docker-compose down
	docker-compose build --no-cache
	docker-compose up -d

install:
	cd backend && go mod download
	npm install

run:
	cd backend && go run cmd/main.go
