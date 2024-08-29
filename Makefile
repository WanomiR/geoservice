include .env

all: clear_logs build

build:
	@docker compose up --build --force-recreate

dump_logs:
	@docker compose logs &> .log

clear_logs:
	@echo "" > .log

swagger_generate: swagger_generate_geo

swagger_generate_geo:
	@cd geo && swag init --parseDependency -g ./cmd/api/main.go && cd ..
