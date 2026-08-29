.PHONY: migrateup migratedown migration test seed gen-docs
MIGRATIONS_PATH=./cmd/migrate/migrations
DB_URL=postgres://admin:admin@localhost/social?sslmode=disable
migration:
	migrate create -seq -ext sql -dir $(MIGRATIONS_PATH) $(filter-out $@,$(MAKECMDGOALS))
migrateup:
	migrate -path=$(MIGRATIONS_PATH) -database "$(DB_URL)" -verbose up
migratedown:
	migrate -path=$(MIGRATIONS_PATH) -database "$(DB_URL)" -verbose down
test:
		go test -v -cover -short ./...
seed:
	go run cmd/migrate/seed/main.go
gen-docs:
	swag init -g ./cmd/api/main.go -o ./docs --parseDependency --parseInternal && swag fmt ./cmd/api/
