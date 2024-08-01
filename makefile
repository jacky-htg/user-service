gen:
	protoc --proto_path=proto proto/users/*.proto --go_out=. --go-grpc_out=. 

init:
	go mod init user-service

migrate:
	go run cmd/cli.go migrate
	
seed:
	go run cmd/cli.go seed

server:
	go run server.go

build:
	env GOOS=linux GOARCH=amd64 go build user-service

.PHONY: gen init migrate seed server build