.PHONY: swagger
swagger:
	@echo "Generando documentación Swagger..."
	@swag init -g cmd/api/main.go -o docs

.PHONY: swagger-install
swagger-install:
	@echo "Instalando Swagger CLI..."
	@go install github.com/swaggo/swag/cmd/swag@v1.8.1

.PHONY: run
run:
	@go run ./cmd/api/main.go

.PHONY: build
build:
	@go build -o app ./cmd/api/*.go
