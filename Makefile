.PHONY: run worker build test up down logs

run:
	go run ./cmd/godance/

worker:
	go run ./cmd/worker/

build:
	go build ./...

test:
	go test ./...

# ── docker-compose ───────────────────────────────────────────────────────────
up:
	docker compose up --build -d

down:
	docker compose down

logs:
	docker compose logs -f app worker
