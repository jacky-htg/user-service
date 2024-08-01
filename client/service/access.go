package service

import (
	"context"
	"log"

	"github.com/jacky-htg/erp-proto/go/pb/users"
)

// ViewAccessTree service client
func ViewAccessTree(ctx context.Context, access users.AccessServiceClient) {
	response, err := access.List(setMetadata(ctx), &users.MyEmpty{})
	if err != nil {
		log.Fatalf("Error when calling grpc service: %s", err)
	}
	log.Printf("Resp received: %v", response)
}
