include .env

all: swagger_generate clear_logs build

build:
	@docker compose up --build --force-recreate

dump_logs:
	@docker compose logs &> .log

clear_logs:
	@echo "" > .log

swagger_generate: swagger_generate_geoservice

swagger_generate_geoservice:
	@cd geoservice && swag init --parseDependency -g ./cmd/api/main.go && cd ..
