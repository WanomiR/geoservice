all: swagger_generate build

build:
	@docker compose up --build --force-recreate

swagger_generate:
	@cd backend && swag init -g ./cmd/api/main.go && cd ..