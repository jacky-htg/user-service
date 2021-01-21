package service

import (
	"context"
	"io"
	"log"
	"user-service/pb/users"
)

// ListFeature service client
func ListFeature(ctx context.Context, feature users.FeatureServiceClient) {
	stream, err := feature.List(setMetadata(ctx), &users.Empty{})
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
		log.Printf("Resp received: %v", resp.GetFeature())
	}
}

// ListPackageFeature service client
func ListPackageFeature(ctx context.Context, packageFeature users.PackageFeatureServiceClient) {
	stream, err := packageFeature.List(setMetadata(ctx), &users.Empty{})
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
		log.Printf("Resp received: %v", resp.GetPackageOfFeature())
	}
}

// ViewPackageFeature service client
func ViewPackageFeature(ctx context.Context, packageFeature users.PackageFeatureServiceClient) {
	response, err := packageFeature.View(setMetadata(ctx), &users.Id{Id: "e1c14424-9ec2-41e1-9709-f65fdaaeddce"})
	if err != nil {
		log.Fatalf("Error when calling grpc service: %s", err)
	}
	log.Printf("Resp received: %v", response)
}
