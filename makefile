all: proto server worker

server:
	go build -o detf_server \
	         cmd/server/*.go

worker:
	go build -o detf_worker \
		 cmd/worker/*.go

proto:
	protoc \
	--go_out=. \
	--go_opt=paths=source_relative \
	--go-grpc_out=. \
	--go-grpc_opt=paths=source_relative \
	api/*.proto

test:
	go test -v ./...
