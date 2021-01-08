package service

import (
	"context"
	"database/sql"
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
