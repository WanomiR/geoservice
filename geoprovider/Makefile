LOCAL_BIN:=$(CURDIR)/bin

generate: generate-rpc-api

install-deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

get-deps:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc

generate-rpc-api:
	mkdir -p ./pkg/geoprovider_v1
	protoc --proto_path ./api/geoprovider_v1 \
	--go_out=./pkg/geoprovider_v1 \
	--go_opt=paths=source_relative \
	--plugin=protoc-gen-go=./bin/protoc-gen-go \
	--go-grpc_out=./pkg/geoprovider_v1 \
	--go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=./bin/protoc-gen-go-grpc \
	./api/geoprovider_v1/geoprovider.proto