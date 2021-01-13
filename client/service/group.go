package service

import (
	"context"
	"io"
	"log"
	users "user-service/pb"
)

// CreateGroup service client
func CreateGroup(ctx context.Context, group users.GroupServiceClient) {
	response, err := group.Create(setMetadata(ctx), &users.Group{
		Name: "Co-Admin",
	})
	if err != nil {
		log.Fatalf("Error when calling grpc service: %s", err)
	}
	log.Printf("Resp received: %v", response)
}

// UpdateGroup service client
func UpdateGroup(ctx context.Context, group users.GroupServiceClient) {
	response, err := group.Update(setMetadata(ctx), &users.Group{
		Id:   "749844e9-d4fa-4121-8a1f-b34c0cba5c02",
		Name: "admin ajah",
	})
	if err != nil {
		log.Fatalf("Error when calling grpc service: %s", err)
	}
	log.Printf("Resp received: %v", response)
}

// ViewGroup service client
func ViewGroup(ctx context.Context, group users.GroupServiceClient) {
	response, err := group.View(setMetadata(ctx), &users.Id{Id: "749844e9-d4fa-4121-8a1f-b34c0cba5c02"})
	if err != nil {
		log.Fatalf("Error when calling grpc service: %s", err)
	}
	log.Printf("Resp received: %v", response)
}

// DeleteGroup service client
func DeleteGroup(ctx context.Context, group users.GroupServiceClient) {
	response, err := group.Delete(setMetadata(ctx), &users.Id{Id: "749844e9-d4fa-4121-8a1f-b34c0cba5c02"})
	if err != nil {
		log.Fatalf("Error when calling grpc service: %s", err)
	}
	log.Printf("Resp received: %v", response)
}

// ListGroup service client
func ListGroup(ctx context.Context, group users.GroupServiceClient) {
	stream, err := group.List(setMetadata(ctx), &users.ListGroupRequest{})
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
		log.Printf("Resp received: %s : %v", resp.GetGroup(), resp.GetPagination())
	}
}

// GrantAccess service client
func GrantAccess(ctx context.Context, group users.GroupServiceClient) {
	response, err := group.GrantAccess(setMetadata(ctx), &users.GrantAccessRequest{
		AccessId: "39306e07-6c0c-4a18-8d02-ae9f63bcad98",
		GroupId:  "ce7ae637-56f8-46c2-ae3c-1be2ed63831d",
	})
	if err != nil {
		log.Fatalf("Error when calling grpc service: %s", err)
	}
	log.Printf("Resp received: %v", response)
}

// RevokeAccess service client
func RevokeAccess(ctx context.Context, group users.GroupServiceClient) {
	response, err := group.RevokeAccess(setMetadata(ctx), &users.GrantAccessRequest{
		AccessId: "39306e07-6c0c-4a18-8d02-ae9f63bcad98",
		GroupId:  "ce7ae637-56f8-46c2-ae3c-1be2ed63831d",
	})
	if err != nil {
		log.Fatalf("Error when calling grpc service: %s", err)
	}
	log.Printf("Resp received: %v", response)
}
