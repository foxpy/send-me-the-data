DB_URL = "postgres://postgres:i_am@localhost:5432/postgres?sslmode=disable"

.PHONY: start-server goose-status

start-server:
	POSTGRES_URL="$(DB_URL)" \
	PREFIX="/home/foxpy/send-me-the-data/dump" \
	go run ./cmd/server

goose-status:
	cd migrations && \
	GOOSE_DRIVER=postgres \
	GOOSE_DBSTRING="$(DB_URL)" \
	goose status

goose-up:
	cd migrations && \
	GOOSE_DRIVER=postgres \
	GOOSE_DBSTRING="$(DB_URL)" \
	goose up

goose-down:
	cd migrations && \
	GOOSE_DRIVER=postgres \
	GOOSE_DBSTRING="$(DB_URL)" \
	goose down-to 0
