protoc:
	rm -f ./pkg/grpc/*.pb.go
	protoc --go_out=. --go_opt=paths=source_relative \
	--go-grpc_out=. --go-grpc_opt=paths=source_relative \
	./pkg/grpc/*.proto

create-migration:
	migrate create -ext sql -dir ./server/internal/adapters/postgres/migrations -seq $(name)

migrate-up:
	migrate -source file://server/internal/adapters/postgres/migrations/ -database postgres://user:password@localhost:5432/gomoney?sslmode=disable up

migrate-down:
	migrate -source file://server/internal/adapters/postgres/migrations/ -database postgres://user:password@localhost:5432/gomoney?sslmode=disable down

docker-rm:
	docker rm gomoney-grpc-server gomoney-api-server gomoney-telegram-server gomoney-db gomoney-cache
	docker rmi gomoney-api gomoney-server-dev gomoney-telegram

sqlc:
	sqlc generate

mock-api:
	mockgen -package mocks -destination ./api/internal/adapters/mock/mock_service.go -source ./api/internal/application/ports.go Service
