include .env

all: swagger_generate build

build:
	@docker compose up --build --force-recreate

swagger_generate:
	@cd geoservice && swag init -g ./cmd/api/main.go && cd ..