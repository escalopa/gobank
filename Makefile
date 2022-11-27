migrateUp:
	migrate -path ./db/migration -database "postgres://postgres:postgres@localhost:5444/bank?sslmode=disable" -verbose up

migrateDown:
	migrate -path ./db/migration -database "postgres://postgres:postgres@localhost:5444/bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/escalopa/go-bank/db/sqlc Store

.PHONY: migrateUp migrateDown sqlc test server mock