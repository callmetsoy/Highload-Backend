.PHONY: proto run-rest run-grpc test bench

# Генерация gRPC кода
proto:
	@mkdir -p gen
	protoc --go_out=gen --go_opt=paths=source_relative \
	       --go-grpc_out=gen --go-grpc_opt=paths=source_relative \
	       proto/urls.proto

# Запуск REST сервера
run-rest:
	go run cmd/rest/main.go

# Запуск gRPC сервера
run-grpc:
	go run cmd/grpc/main.go

# Тесты
test:
	go test -v ./...

# Бенчмарки (нужно сначала запустить gRPC сервер)
bench:
	go test -v -bench=. -benchmem ./benchmark/