all: proto server worker

server:
	go build cmd/server/detf_server.go

worker:
	go build cmd/worker/detf_worker.go

proto:
	protoc \
	--go_out=. \
	--go_opt=paths=source_relative \
	--go-grpc_out=. \
	--go-grpc_opt=paths=source_relative \
	api/*.proto
