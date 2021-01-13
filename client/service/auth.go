package service

import (
	"context"
	"log"
	"user-service/internal/pkg/app"
	users "user-service/pb"
)

// Login service client
func Login(ctx context.Context, auth users.AuthServiceClient) context.Context {
	response, err := auth.Login(ctx, &users.LoginRequest{Username: "wira-admin", Password: "1234"})
	if err != nil {
		log.Fatalf("Error when calling grpc service: %s", err)
	}
	log.Printf("Resp received: %v", response)

	ctx = context.WithValue(ctx, app.Ctx("token"), response.GetToken())
	ctx = context.WithValue(ctx, app.Ctx("userID"), response.GetUser().GetId())
	ctx = context.WithValue(ctx, app.Ctx("companyID"), response.GetUser().GetCompanyId())

	return ctx
}

// ForgotPassword service client
func ForgotPassword(ctx context.Context, auth users.AuthServiceClient) {
	response, err := auth.ForgotPassword(ctx, &users.ForgotPasswordRequest{Email: "rijal.asep.nugroho@gmail.com"})
	if err != nil {
		log.Fatalf("Error when calling grpc service: %s", err)
	}
	log.Printf("Resp received: %v", response)
}

// ResetPassword service client
func ResetPassword(ctx context.Context, auth users.AuthServiceClient) {
	response, err := auth.ResetPassword(ctx, &users.ResetPasswordRequest{
		Token:       "4d5de8c6-46cd-453f-88f1-fe904ee01746",
		NewPassword: "12345",
		RePassword:  "12345",
	})

	if err != nil {
		log.Fatalf("Error when calling grpc service: %s", err)
	}
	log.Printf("Resp received: %v", response)
}

// ChangePassword service client
func ChangePassword(ctx context.Context, auth users.AuthServiceClient) {
	response, err := auth.ChangePassword(setMetadata(ctx), &users.ChangePasswordRequest{
		OldPassword: "12345",
		NewPassword: "1234",
		RePassword:  "1234",
	})

	if err != nil {
		log.Fatalf("Error when calling grpc service: %s", err)
	}
	log.Printf("Resp received: %v", response)
}

// IsAuth service client
func IsAuth(ctx context.Context, auth users.AuthServiceClient) {
	response, err := auth.IsAuth(setMetadata(ctx), &users.String{String_: "asal"})

	if err != nil {
		log.Fatalf("Error when calling grpc service: %s", err)
	}
	log.Printf("Resp received: %v", response)
}
