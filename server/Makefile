include .env

CURRENT_DIR = $(shell pwd)

DB_URL := postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

mig-up:
	migrate -path internal/migrations -database '${DB_URL}' -verbose up

mig-down:
	migrate -path internal/migrations -database '${DB_URL}' -verbose down

mig-force:
	migrate -path internal/migrations -database '${DB_URL}' -verbose force 1

mig-create-tables:
	migrate create -ext sql -dir internal/migrations -seq DB_NAME

swag-gen:
	~/go/bin/swag init -g cmd/sentinel/main.go -o internal/handlers/docs

print:
	@echo $(DB_USER)