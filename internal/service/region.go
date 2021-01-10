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
	var output users.Region
	var err error
	var regionModel model.Region

	// basic validation
	{
		if len(in.GetId()) == 0 {
			return &output, status.Error(codes.InvalidArgument, "Please supply valid id")
		}
		regionModel.Pb.Id = in.GetId()
	}

	ctx, err = getMetadata(ctx)
	if err != nil {
		return &output, err
	}

	// get user login
	var userLogin model.User
	userLogin.Pb.Id = ctx.Value(app.Ctx("userID")).(string)
	err = userLogin.Get(ctx, u.Db)
	if err != nil {
		return &output, err
	}

	if len(userLogin.Pb.GetBranchId()) > 0 {
		return &output, status.Error(codes.Unauthenticated, "only user company/region can update the region")
	}

	if len(userLogin.Pb.GetRegionId()) > 0 && userLogin.Pb.GetRegionId() != in.GetId() {
		return &output, status.Error(codes.Unauthenticated, "its not your region")
	}

	err = regionModel.Get(ctx, u.Db)
	if err != nil {
		return &output, err
	}

	if userLogin.Pb.GetCompanyId() != regionModel.Pb.GetCompanyId() {
		return &output, status.Error(codes.Unauthenticated, "its not your company")
	}

	if len(in.GetName()) > 0 {
		regionModel.Pb.Name = in.GetName()
	}

	if len(in.GetBranches()) > 0 {
		regionModel.UpdateBranches = true
		for _, branch := range in.GetBranches() {
			var branchModel model.Branch
			branchModel.Pb.Id = branch.GetId()
			err = branchModel.Get(ctx, u.Db)
			if err != nil {
				return &output, err
			}
		}

		regionModel.Pb.Branches = in.GetBranches()
	}

	tx, err := u.Db.BeginTx(ctx, nil)
	if err != nil {
		return &output, status.Errorf(codes.Internal, "begin tx: %v", err)
	}

	err = regionModel.Update(ctx, u.Db, tx)
	if err != nil {
		tx.Rollback()
		return &output, err
	}

	tx.Commit()

	return &regionModel.Pb, nil
}

// View Region
func (u *Region) View(ctx context.Context, in *users.Id) (*users.Region, error) {
	var output users.Region
	var err error
	var regionModel model.Region

	// basic validation
	{
		if len(in.GetId()) == 0 {
			return &output, status.Error(codes.InvalidArgument, "Please supply valid id")
		}
		regionModel.Pb.Id = in.GetId()
	}

	ctx, err = getMetadata(ctx)
	if err != nil {
		return &output, err
	}

	// get user login
	var userLogin model.User
	userLogin.Pb.Id = ctx.Value(app.Ctx("userID")).(string)
	err = userLogin.Get(ctx, u.Db)
	if err != nil {
		return &output, err
	}

	err = regionModel.Get(ctx, u.Db)
	if err != nil {
		return &output, err
	}

	if userLogin.Pb.GetCompanyId() != regionModel.Pb.GetCompanyId() {
		return &output, status.Error(codes.Unauthenticated, "its not your company")
	}

	if len(userLogin.Pb.GetRegionId()) > 0 && userLogin.Pb.GetRegionId() != regionModel.Pb.GetId() {
		return &output, status.Error(codes.Unauthenticated, "its not your region")
	}

	return &regionModel.Pb, nil
}

// Delete Region
func (u *Region) Delete(ctx context.Context, in *users.Id) (*users.Boolean, error) {
	var output users.Boolean
	output.Boolean = false

	var err error
	var regionModel model.Region

	// basic validation
	{
		if len(in.GetId()) == 0 {
			return &output, status.Error(codes.InvalidArgument, "Please supply valid id")
		}
		regionModel.Pb.Id = in.GetId()
	}

	ctx, err = getMetadata(ctx)
	if err != nil {
		return &output, err
	}

	// get user login
	var userLogin model.User
	userLogin.Pb.Id = ctx.Value(app.Ctx("userID")).(string)
	err = userLogin.Get(ctx, u.Db)
	if err != nil {
		return &output, err
	}

	if len(userLogin.Pb.GetRegionId()) > 0 || len(userLogin.Pb.GetBranchId()) > 0 {
		return &output, status.Error(codes.Unauthenticated, "only user company can delete region")
	}

	err = regionModel.Get(ctx, u.Db)
	if err != nil {
		return &output, err
	}

	if userLogin.Pb.GetCompanyId() != regionModel.Pb.GetCompanyId() {
		return &output, status.Error(codes.Unauthenticated, "its not your company")
	}

	err = regionModel.Delete(ctx, u.Db)
	if err != nil {
		return &output, err
	}

	output.Boolean = true
	return &output, nil
}

// List Region
func (u *Region) List(ctx context.Context, in *users.ListRegionRequest) (*users.ListRegionResponse, error) {
	return &users.ListRegionResponse{}, nil
}
