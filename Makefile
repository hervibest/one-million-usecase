DB_URL=postgres://root:password@localhost:5433/one_million_usecase?sslmode=disable&TimeZone=Asia/Jakarta
MIGRATION_DIR=db/migrations

migrate-up:
	goose -dir $(MIGRATION_DIR) postgres "$(DB_URL)" up

migrate-down:
	goose -dir $(MIGRATION_DIR) postgres "$(DB_URL)" down

migrate-status:
	goose -dir $(MIGRATION_DIR) postgres "$(DB_URL)" status

migrate-create:
	goose -dir $(MIGRATION_DIR) create $(name) sql

.PHONY: migrate-up migrate-down migrate-status migrate-create