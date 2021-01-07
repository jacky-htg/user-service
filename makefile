gen:
	protoc --proto_path=proto proto/*.proto --go_out=plugins=grpc:. 

init:
	go mod init user-service

server:
	go run server.go

.PHONY: gen init server