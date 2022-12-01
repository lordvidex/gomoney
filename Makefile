protoc: 
	protoc --go_out=. --go_opt=paths=source_relative \
	--go-grpc_out=. --go-grpc_opt=paths=source_relative \
	./server/internal/adapters/grpc/*.proto

create-migration:
	migrate create -ext sql -dir ./server/internal/adapters/postgres/migrations -seq $(name)

migrate-up:
	migrate -source file://server/internal/adapters/postgres/migrations/ -database postgres://user:password@localhost:5432/gomoney?sslmode=disable up
migrate-down:
	migrate -source file://server/internal/adapters/postgres/migrations/ -database postgres://user:password@localhost:5432/gomoney?sslmode=disable down

sqlc:
	sqlc generate