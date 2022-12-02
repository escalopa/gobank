migrateUp:
	migrate -path ./db/migration -database "postgres://postgres:postgres@localhost:5444/bank?sslmode=disable" -verbose up $(version)

migrateDown:
	migrate -path ./db/migration -database "postgres://postgres:postgres@localhost:5444/bank?sslmode=disable" -verbose down $(version)

migrateCreate:
	migrate create -ext sql -dir ./db/migration -seq $(name)

migrateForce:
	migrate -path ./db/migration -database "postgres://postgres:postgres@localhost:5444/bank?sslmode=disable" -verbose force $(version)

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/escalopa/gobank/db/sqlc Store

proto:
	rm -rf pb/*
	protoc  --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
    proto/*.proto

evans: 
	evans --host localhost --port 9000 -r repl

.PHONY: migrateUp migrateDown sqlc test server mock migrateCreate proto evans