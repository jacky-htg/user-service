package service

import (
	"context"
	"io"
	"log"
	"user-service/pb/users"
)

// CreateRegion service client
func CreateRegion(ctx context.Context, region users.RegionServiceClient) {
	response, err := region.Create(setMetadata(ctx), &users.Region{
		Code: "WEST",
		Name: "Sumatera",
	})
	if err != nil {
		log.Fatalf("Error when calling grpc service: %s", err)
	}
	log.Printf("Resp received: %v", response)
}

// UpdateRegion service client
func UpdateRegion(ctx context.Context, region users.RegionServiceClient) {
	response, err := region.Update(setMetadata(ctx), &users.Region{
		Id:   "872df9ca-8758-4a50-a49b-b9ebc41314c5",
		Name: "Sumatera Raya",
	})
	if err != nil {
		log.Fatalf("Error when calling grpc service: %s", err)
	}
	log.Printf("Resp received: %v", response)
}

// ViewRegion service client
func ViewRegion(ctx context.Context, region users.RegionServiceClient) {
	response, err := region.View(setMetadata(ctx), &users.Id{Id: "872df9ca-8758-4a50-a49b-b9ebc41314c5"})
	if err != nil {
		log.Fatalf("Error when calling grpc service: %s", err)
	}
	log.Printf("Resp received: %v", response)
}

// DeleteRegion service client
func DeleteRegion(ctx context.Context, region users.RegionServiceClient) {
	response, err := region.Delete(setMetadata(ctx), &users.Id{Id: "872df9ca-8758-4a50-a49b-b9ebc41314c5"})
	if err != nil {
		log.Fatalf("Error when calling grpc service: %s", err)
	}
	log.Printf("Resp received: %v", response)
}

// ListRegion service client
func ListRegion(ctx context.Context, region users.RegionServiceClient) {
	stream, err := region.List(setMetadata(ctx), &users.ListRegionRequest{})
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
		log.Printf("Resp received: %s : %v", resp.GetRegion(), resp.GetPagination())
	}
}
