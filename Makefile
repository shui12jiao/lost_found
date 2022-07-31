postgres:
	docker run --name postgres14 --network lost_found -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:14-alpine
	
createdb:
	docker exec -it postgres14 createdb --username=root --owner=root lost_found
	
dropdb:
	docker exec -it postgres14 dropdb lost_found

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@:5432/lost_found?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/lost_found?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb  -destination db/mock/store.go github.com/shui12jiao/lost_found/db/sqlc Store

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test server mock
