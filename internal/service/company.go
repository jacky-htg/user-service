package service

import (
	"context"
	"database/sql"
	"user-service/internal/pkg/db/redis"
	users "user-service/pb"
)

// Company struct
type Company struct {
	Db    *sql.DB
	Cache *redis.Cache
}

// Register new company
func (u *Company) Register(ctx context.Context, in *users.CompanyRegister) (*users.CompanyRegister, error) {
	var output users.CompanyRegister
	var err error

	return &output, err
}

// Update Company
func (u *Company) Update(ctx context.Context, in *users.Company) (*users.Company, error) {
	var output users.Company
	var err error

	return &output, err
}

// View Company
func (u *Company) View(ctx context.Context, in *users.Id) (*users.Company, error) {
	var output users.Company
	var err error

	return &output, err
}

// Delete Company
func (u *Company) Delete(ctx context.Context, in *users.Id) (*users.Boolean, error) {
	var output users.Boolean
	var err error

	return &output, err
}

// List Companies
func (u *Company) List(in *users.ListCompanyRequest, stream users.CompanyService_ListServer) error {
	var err error
	return err
}
