package service

import (
	"context"
	"io"
	"log"
	users "user-service/pb"
)

// CreateBranch service client
func CreateBranch(ctx context.Context, branch users.BranchServiceClient) {
	response, err := branch.Create(setMetadata(ctx), &users.Branch{
		Code:     "LAPG",
		Name:     "Lampung Branch",
		Address:  "Jl Address",
		City:     "Lampung",
		Province: "Lampung",
		Phone:    "0814111111111",
		Pic:      "Mr Lampung",
		PicPhone: "081122222222",
		RegionId: "69927081-7fd4-4c2f-9ac1-32647aa055a7",
	})
	if err != nil {
		log.Fatalf("Error when calling grpc service: %s", err)
	}
	log.Printf("Resp received: %v", response)
}

// UpdateBranch service client
func UpdateBranch(ctx context.Context, branch users.BranchServiceClient) {
	response, err := branch.Update(setMetadata(ctx), &users.Branch{
		Id:       "f174f588-c10d-4ac4-ade3-6ab326764bf5",
		Name:     "Lampung Branch",
		Address:  "Jl Address",
		City:     "Lampung",
		Province: "Lampung",
		Phone:    "0814111111111",
		Pic:      "Mr Lampung",
		PicPhone: "081122222222",
		RegionId: "69927081-7fd4-4c2f-9ac1-32647aa055a7",
	})
	if err != nil {
		log.Fatalf("Error when calling grpc service: %s", err)
	}
	log.Printf("Resp received: %v", response)
}

// ViewBranch service client
func ViewBranch(ctx context.Context, branch users.BranchServiceClient) {
	response, err := branch.View(setMetadata(ctx), &users.Id{Id: "f174f588-c10d-4ac4-ade3-6ab326764bf5"})
	if err != nil {
		log.Fatalf("Error when calling grpc service: %s", err)
	}
	log.Printf("Resp received: %v", response)
}

// DeleteBranch service client
func DeleteBranch(ctx context.Context, branch users.BranchServiceClient) {
	response, err := branch.Delete(setMetadata(ctx), &users.Id{Id: "f174f588-c10d-4ac4-ade3-6ab326764bf5"})
	if err != nil {
		log.Fatalf("Error when calling grpc service: %s", err)
	}
	log.Printf("Resp received: %v", response)
}

// ListBranch service client
func ListBranch(ctx context.Context, branch users.BranchServiceClient) {
	stream, err := branch.List(setMetadata(ctx), &users.ListBranchRequest{})
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
		log.Printf("Resp received: %s : %v", resp.GetBranch(), resp.GetPagination())
	}
}
