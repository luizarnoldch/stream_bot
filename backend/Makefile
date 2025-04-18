DB_USER       := postgres
DB_PASS       := postgres
DB_HOST       := localhost
DB_PORT       := 5434
DB_NAME       := postgres
SSL_MODE      := disable
DB_URL        := postgres://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(SSL_MODE)

MIGRATIONS_DIR := db/sql/schemas

.PHONY: ng sql-up sql-down test

ng:
	./ngrok-local.sh

sql-up:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" up

sql-down:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" down
gen:
	cd db/sql && sqlc generate