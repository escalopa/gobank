migrateUp:
	migrate -path ./db/migration -database "postgres://postgres:postgres@localhost:5444/bank?sslmode=disable" -verbose up $(version)

migrateDown:
	migrate -path ./db/migration -database "postgres://postgres:postgres@localhost:5444/bank?sslmode=disable" -verbose down $(version)

migrateCreate:
	migrate create -ext sql -dir ./db/migration -seq $(name)

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/escalopa/gobank/db/sqlc Store

.PHONY: migrateUp migrateDown sqlc test server mock migrateCreate