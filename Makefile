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

gendoc: 
	sql2dbml --postgres ./dcos/schema.sql -o docs/db.dbml

proto:
	rm -rf pb/*
	rm -rf docs/swagger/*
	protoc  --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
		--grpc-gateway_out pb --grpc-gateway_opt=logtostderr=true --grpc-gateway_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
 		--openapiv2_out=docs/swagger/ --openapiv2_opt=allow_merge=true,merge_file_name=gobank \
    proto/*.proto

evans: 
	evans --host localhost --port 9000 -r repl

.PHONY: migrateUp migrateDown sqlc test server mock migrateCreate proto evans gendocs