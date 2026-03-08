DB_URL = "postgres://postgres:i_am@localhost:5432/postgres?sslmode=disable"
MIGRATIONS = "cmd/server/idb/postgres/migrations/"

.PHONY: start-server goose-status

start-server:
	POSTGRES_URL="$(DB_URL)" \
	PREFIX="/home/foxpy/send-me-the-data/dump" \
	go run ./cmd/server

goose-status:
	GOOSE_DRIVER=postgres \
	GOOSE_DBSTRING=$(DB_URL) \
	GOOSE_MIGRATION_DIR=$(MIGRATIONS) \
	goose status

goose-up:
	GOOSE_DRIVER=postgres \
	GOOSE_DBSTRING=$(DB_URL) \
	GOOSE_MIGRATION_DIR=$(MIGRATIONS) \
	goose up

goose-down:
	GOOSE_DRIVER=postgres \
	GOOSE_DBSTRING=$(DB_URL) \
	GOOSE_MIGRATION_DIR=$(MIGRATIONS) \
	goose down-to 0
