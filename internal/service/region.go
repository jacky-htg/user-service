package service

import (
	"context"
	"database/sql"
	"user-service/internal/model"
	"user-service/internal/pkg/app"
	"user-service/internal/pkg/db/redis"
	users "user-service/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Region struct
type Region struct {
	Db    *sql.DB
	Cache *redis.Cache
}

// Create new region
func (u *Region) Create(ctx context.Context, in *users.Region) (*users.Region, error) {
	var output users.Region
	var err error
	var regionModel model.Region

	// basic validation
	{
		if len(in.GetName()) == 0 {
			return &output, status.Error(codes.InvalidArgument, "Please supply valid name")
		}
	}

	// code validation
	{
		if len(in.GetCode()) == 0 {
			return &output, status.Error(codes.InvalidArgument, "Please supply valid code")
		}

		regionModel = model.Region{}
		regionModel.Pb.Code = in.GetCode()
		err = regionModel.GetByCode(ctx, u.Db)
		if err != nil {
			if st, ok := status.FromError(err); ok && st.Code() != codes.NotFound {
				return &output, err
			}
		}

		if len(regionModel.Pb.GetId()) > 0 {
			return &output, status.Error(codes.AlreadyExists, "code must be unique")
		}
	}

	ctx, err = getMetadata(ctx)
	if err != nil {
		return &output, err
	}

	// company validation
	{
		if len(in.GetCompanyId()) > 0 && in.GetCompanyId() != ctx.Value(app.Ctx("companyID")).(string) {
			return &output, status.Error(codes.PermissionDenied, "Please supply valid company id")
		}
		in.CompanyId = ctx.Value(app.Ctx("companyID")).(string)
	}

	// get user login
	var userLogin model.User
	userLogin.Pb.Id = ctx.Value(app.Ctx("userID")).(string)
	err = userLogin.Get(ctx, u.Db)
	if err != nil {
		return &output, err
	}

	if len(userLogin.Pb.GetRegionId()) > 0 || len(userLogin.Pb.GetBranchId()) > 0 {
		return &output, status.Error(codes.PermissionDenied, "Need user company to add new region")
	}

	if len(in.GetBranches()) > 0 {
		for _, branch := range in.GetBranches() {
			branchModel := model.Branch{}
			branchModel.Pb.Id = branch.GetId()
			err = branchModel.Get(ctx, u.Db)
			if err != nil {
				return &output, err
			}
		}
	}

	regionModel = model.Region{}
	regionModel.Pb = users.Region{
		Branches:  in.GetBranches(),
		Code:      in.GetCode(),
		CompanyId: in.GetCompanyId(),
		Name:      in.GetName(),
	}
	tx, err := u.Db.BeginTx(ctx, nil)
	if err != nil {
		return &output, status.Errorf(codes.Internal, "begin tx: %v", err)
	}

	err = regionModel.Create(ctx, u.Db, tx)
	if err != nil {
		tx.Rollback()
		return &output, err
	}

	tx.Commit()

	return &regionModel.Pb, nil
}

// Update region
func (u *Region) Update(ctx context.Context, in *users.Region) (*users.Region, error) {
	return &users.Region{}, nil
}

// View Region
func (u *Region) View(ctx context.Context, in *users.Id) (*users.Region, error) {
	return &users.Region{}, nil
}

// Delete Region
func (u *Region) Delete(ctx context.Context, in *users.Id) (*users.Boolean, error) {
	return &users.Boolean{}, nil
}

// List Region
func (u *Region) List(ctx context.Context, in *users.ListRegionRequest) (*users.ListRegionResponse, error) {
	return &users.ListRegionResponse{}, nil
}
