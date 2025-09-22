init:
	go mod init github.com/Viczdera/bank
	
createdb:
	docker exec -it bank_postgres12 createdb --username=root --owner=root s_bank

dropdb:
	docker exec -it bank_postgres12 dropdb s_bank

postgres:
	docker run --name bank_postgres12 --network bank-network -p 8080:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

migrateup:
	migrate -path db/migrations  -database "postgres://root:secret@localhost:8080/s_bank?sslmode=disable" -verbose up

migrateup1:
	migrate -path db/migrations  -database "postgres://root:secret@localhost:8080/s_bank?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migrations  -database "postgres://root:secret@localhost:8080/s_bank?sslmode=disable" -verbose down

migratedown1:
	migrate -path db/migrations  -database "postgres://root:secret@localhost:8080/s_bank?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/Viczdera/bank-gorm/db/sqlc Store

.PHONY:
	postgres createdb dropdb migrateup migratedown migrateup1 migratedown1 test server mock