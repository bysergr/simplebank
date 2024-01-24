createPostgresContainer:
	docker run --name postgres16 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:16-alpine

dropPostgresContainer:
	docker rm postgres16

createDB:
	docker exec -it postgres16 createdb --username=root --owner=root simple_bank

dropDB:
	docker exec -it postgres16 dropdb simple_bank

migrateUp:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

migrateDown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -destination db/mock/store.go -package mockdb github.com/bysergr/simple-bank/db/sqlc Store


.PHONY: createPostgresContainer dropPostgresContainer createDB dropDB migrateUp migrateDown sqlc test server mock
