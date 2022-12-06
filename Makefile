protoc:
	rm -f ./pkg/grpc/*.pb.go
	protoc --experimental_allow_proto3_optional --go_out=. --go_opt=paths=source_relative \
	--go-grpc_out=. --go-grpc_opt=paths=source_relative \
	./pkg/grpc/*.proto

create-migration:
	migrate create -ext sql -dir ./server/internal/adapters/postgres/migrations -seq $(name)

migrate-up:
	migrate -source file://server/internal/adapters/postgres/migrations/ -database postgres://user:password@localhost:5432/gomoney?sslmode=disable up

migrate-down:
	migrate -source file://server/internal/adapters/postgres/migrations/ -database postgres://user:password@localhost:5432/gomoney?sslmode=disable down

docker-rm:
	docker rm gomoney-grpc-server gomoney-api-server gomoney-telegram-server gomoney-db
	docker rmi gomoney-api gomoney-server-dev gomoney-telegram

sqlc:
	sqlc generate