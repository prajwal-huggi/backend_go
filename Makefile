include .env

DB_URL=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

migrate-up:
	goose -dir migrations postgres "$(DB_URL)" up

migrate-down:
	goose -dir migrations postgres "$(DB_URL)" down

migrate-reset:
	goose -dir migrations postgres "$(DB_URL)" reset

migrate-status:
	goose -dir migrations postgres "$(DB_URL)" status

# COMMANDS TO RUN 
# make migrate-up, migrate-down, migrate-reset