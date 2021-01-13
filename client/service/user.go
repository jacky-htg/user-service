package service

import (
	"io"
	"log"
	users "user-service/pb"

	"golang.org/x/net/context"
)

// CreateUser service Client
func CreateUser(ctx context.Context, user users.UserServiceClient) {
	response, err := user.Create(setMetadata(ctx), &users.User{
		Email:    "admin@gmail.com",
		Group:    &users.Group{Id: "b7365e19-4f39-43bd-9b4a-f53175ab37bd"},
		Name:     "Administrator",
		Username: "jackyhtg",
	})

	if err != nil {
		log.Fatalf("Error when calling grpc service: %s", err)
	}
	log.Printf("Resp received: %v", response)
}

// UpdateUser service client
func UpdateUser(ctx context.Context, user users.UserServiceClient) {
	response, err := user.Create(setMetadata(ctx), &users.User{
		Email:    "admin@gmail.com",
		Group:    &users.Group{Id: "b7365e19-4f39-43bd-9b4a-f53175ab37bd"},
		Name:     "Administrator",
		Username: "jackyhtg",
	})

	if err != nil {
		log.Fatalf("Error when calling grpc service: %s", err)
	}
	log.Printf("Resp received: %v", response)
}

// ViewUser service cient
func ViewUser(ctx context.Context, user users.UserServiceClient) {
	response, err := user.View(setMetadata(ctx), &users.Id{Id: "f65b47d0-4fb6-4418-b97b-d736106c857e"})

	if err != nil {
		log.Fatalf("Error when calling grpc service: %s", err)
	}
	log.Printf("Resp received: %v", response)
}

// DeleteUser service client
func DeleteUser(ctx context.Context, user users.UserServiceClient) {
	response, err := user.Delete(setMetadata(ctx), &users.Id{Id: "f65b47d0-4fb6-4418-b97b-d736106c857e"})

	if err != nil {
		log.Fatalf("Error when calling grpc service: %s", err)
	}
	log.Printf("Resp received: %v", response)
}

// ListUser service client
func ListUser(ctx context.Context, user users.UserServiceClient) {
	stream, err := user.List(setMetadata(ctx), &users.ListUserRequest{})
	if err != nil {
		log.Fatalf("Error when calling grpc service: %s", err)
	}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			log.Fatal("end stream")
			break
		}
		if err != nil {
			log.Fatalf("cannot receive %v", err)
		}
		log.Printf("Resp received: %s : %v", resp.GetUser(), resp.GetPagination())
	}
}

// GetUserByToken service client
func GetUserByToken(ctx context.Context, user users.UserServiceClient) {
	response, err := user.GetByToken(setMetadataToken(ctx), &users.Empty{})

	if err != nil {
		log.Fatalf("Error when calling grpc service: %s", err)
	}
	log.Printf("Resp received: %v", response)
}
