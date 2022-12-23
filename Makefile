migrateUp:
	migrate -path ./db/migration -database "postgres://postgres:postgres@localhost:5432/gobank?sslmode=disable" -verbose up $(version)

migrateDown:
	migrate -path ./db/migration -database "postgres://postgres:postgres@localhost:5432/gobank?sslmode=disable" -verbose down $(version)

migrateCreate:
	migrate create -ext sql -dir ./db/migration -seq $(name)

migrateForce:
	migrate -path ./db/migration -database "postgres://postgres:postgres@localhost:5432/gobank?sslmode=disable" -verbose force $(version)

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/escalopa/gobank/db/sqlc Store

gendoc: 
	sql2dbml --postgres ./dcos/schema.sql -o docs/db.dbml

proto:
	rm -rf grpc/pb/*
	rm -rf docs/swagger/json/*
	protoc  --proto_path=grpc/proto --go_out=grpc/pb --go_opt=paths=source_relative \
		--grpc-gateway_out grpc/pb --grpc-gateway_opt=logtostderr=true --grpc-gateway_opt=paths=source_relative \
    --go-grpc_out=grpc/pb --go-grpc_opt=paths=source_relative \
 		--openapiv2_out=docs/swagger/json/ --openapiv2_opt=allow_merge=true,merge_file_name=gobank \
		--experimental_allow_proto3_optional \
    grpc/proto/*.proto

evans: 
	evans --host localhost --port 9000 -r repl

swagger:
	swag fmt & swag init -d ./api/handlers/ -g ../cmd/main.go -o ./api/docs/ --parseDependency

.PHONY: migrateUp migrateDown sqlc test server mock migrateCreate proto evans gendocs swagger