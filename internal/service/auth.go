package service

import (
	"context"
	"database/sql"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"user-service/internal/pkg/app"
	"user-service/internal/pkg/db/redis"
	users "user-service/pb"
)

// Auth struct
type Auth struct {
	Db    *sql.DB
	Cache *redis.Cache
}

// Login service
func (u *Auth) Login(context.Context, *users.LoginRequest) (*users.LoginResponse, error) {
	return &users.LoginResponse{}, nil
}

// ForgotPassword service
func (u *Auth) ForgotPassword(context.Context, *users.ForgotPasswordRequest) (*users.Message, error) {
	return &users.Message{}, nil
}

// ResetPassword service
func (u *Auth) ResetPassword(context.Context, *users.ResetPasswordRequest) (*users.Message, error) {
	return &users.Message{}, nil
}

// ChangePassword service
func (u *Auth) ChangePassword(context.Context, *users.ChangePasswordRequest) (*users.Message, error) {
	return &users.Message{}, nil
}

// IsAuth service
func (u *Auth) IsAuth(context.Context, *users.Id) (*users.Boolean, error) {
	return &users.Boolean{}, nil
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
