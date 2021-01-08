package service

import (
	"context"
	"database/sql"
	"user-service/internal/model"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"user-service/internal/pkg/app"
	"user-service/internal/pkg/db/redis"
	"user-service/internal/pkg/token"
	users "user-service/pb"
)

// Auth struct
type Auth struct {
	Db    *sql.DB
	Cache *redis.Cache
}

// Login service
func (u *Auth) Login(ctx context.Context, in *users.LoginRequest) (*users.LoginResponse, error) {
	var output users.LoginResponse
	if len(in.GetUsername()) <= 0 {
		return &output, status.Error(codes.InvalidArgument, "Please supply valid username")
	}

	if len(in.GetPassword()) <= 0 {
		return &output, status.Error(codes.InvalidArgument, "Please supply valid password")
	}

	var userModel model.User
	userModel.Pb.Username = in.GetUsername()
	err := userModel.GetByUserNamePassword(ctx, u.Db, in.GetPassword())
	if err != nil {
		return &output, err
	}

	output.User = &userModel.Pb
	output.Token, err = token.ClaimToken(output.User.GetEmail())
	if err != nil {
		return &output, status.Errorf(codes.Internal, "claim token: %v", err)
	}

	return &output, nil
}

// ForgotPassword service
func (u *Auth) ForgotPassword(ctx context.Context, in *users.ForgotPasswordRequest) (*users.Message, error) {
	return &users.Message{}, nil
}

// ResetPassword service
func (u *Auth) ResetPassword(ctx context.Context, in *users.ResetPasswordRequest) (*users.Message, error) {
	return &users.Message{}, nil
}

// ChangePassword service
func (u *Auth) ChangePassword(ctx context.Context, in *users.ChangePasswordRequest) (*users.Message, error) {
	return &users.Message{}, nil
}

// IsAuth service
func (u *Auth) IsAuth(ctx context.Context, in *users.Id) (*users.Boolean, error) {
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
