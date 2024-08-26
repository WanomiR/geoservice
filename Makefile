include .env

all: swagger_generate build

build:
	@docker compose up --build --force-recreate

swagger_generate: swagger_generate_geoservice

swagger_generate_geoservice:
	@cd geoservice && swag init --parseDependency -g ./cmd/api/main.go && cd ..
