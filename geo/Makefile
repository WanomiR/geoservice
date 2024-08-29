LOCAL_BIN:=$(CURDIR)/bin
PACKAGE_NAME:=geo_v1
SERVICE_NAME:=geo

generate: generate-grpc-api

install-deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

get-deps:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc

generate-grpc-api:
	mkdir -p ./pkg/$(PACKAGE_NAME)
	protoc --proto_path ./api/$(PACKAGE_NAME) \
	--go_out=./pkg/$(PACKAGE_NAME) \
	--go_opt=paths=source_relative \
	--plugin=protoc-gen-go=./bin/protoc-gen-go \
	--go-grpc_out=./pkg/$(PACKAGE_NAME) \
	--go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=./bin/protoc-gen-go-grpc \
	./api/$(PACKAGE_NAME)/$(SERVICE_NAME).proto