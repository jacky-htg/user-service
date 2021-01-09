package service

import (
	"context"
	"user-service/internal/pkg/app"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func contextError(ctx context.Context) error {
	switch ctx.Err() {
	case context.Canceled:
		return status.Error(codes.Canceled, "request is canceled")
	case context.DeadlineExceeded:
		return status.Error(codes.DeadlineExceeded, "deadline is exceeded")
	default:
		return nil
	}
}

func getMetadataToken(ctx context.Context) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	token := md["token"]
	if len(token) == 0 {
		return ctx, status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}

	ctx = context.WithValue(ctx, app.Ctx("token"), token[0])

	return ctx, nil
}

func getMetadata(ctx context.Context) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	userID := md["user_id"]
	if len(userID) == 0 {
		return ctx, status.Errorf(codes.Unauthenticated, "user_id is not provided")
	}

	ctx = context.WithValue(ctx, app.Ctx("userID"), userID[0])

	companyID := md["company_id"]
	if len(companyID) == 0 {
		return ctx, status.Errorf(codes.Unauthenticated, "company_id is not provided")
	}

	ctx = context.WithValue(ctx, app.Ctx("companyID"), companyID[0])

	return ctx, nil
}
