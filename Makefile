-include .env

DB_URL=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

build:
	go build ./...

test:
	go test ./... -v

vet:
	go vet ./...

fmt:
	gofmt -w .

migrate-up:
	goose -dir migrations postgres "$(DB_URL)" up

migrate-down:
	goose -dir migrations postgres "$(DB_URL)" down

migrate-reset:
	goose -dir migrations postgres "$(DB_URL)" reset

migrate-status:
	goose -dir migrations postgres "$(DB_URL)" status

ci:
	make fmt
	make build
	make vet
	make migrate-up
	make test
	docker build -t backend-go .
# COMMANDS TO RUN 
# make migrate-up, migrate-down, migrate-reset