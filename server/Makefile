CURRENT_DIR = $(shell pwd)

DB_URL := postgres://postgres:03212164@localhost:5432/database_name?sslmode=disable

swag-init:
	swag init -g api/router.go --output api/handler/docs


mig-up:
	migrate -path db/migrations -database '${DB_URL}' -verbose up

mig-down:
	migrate -path db/migrations -database '${DB_URL}' -verbose down

mig-force:
	migrate -path db/migrations -database '${DB_URL}' -verbose force 1

mig-create-tables:
	migrate create -ext sql -dir internal/migrations -seq monitoring_db

swag-gen:
	~/go/bin/swag init -g api/router.go -o api/handler/docs