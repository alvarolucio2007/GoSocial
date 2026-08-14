.PHONY: migrateup migratedown migration
MIGRATIONS_PATH=./cmd/migrate/migrations
DB_URL=postgres://admin:admin@localhost/social?sslmode=disable
migration:
	migrate create -seq -ext sql -dir $(MIGRATIONS_PATH) $(filter-out $@,$(MAKECMDGOALS))
migrateup:
	migrate -path=$(MIGRATIONS_PATH) -database "$(DB_URL)" -verbose up
migratedown:
	migrate -path=$(MIGRATIONS_PATH) -database "$(DB_URL)" -verbose down
