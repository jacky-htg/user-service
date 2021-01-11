package service

import (
	"context"
	"database/sql"
	users "user-service/pb"

	"user-service/internal/model"
	"user-service/internal/pkg/app"
	"user-service/internal/pkg/db/redis"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Branch struct
type Branch struct {
	Db    *sql.DB
	Cache *redis.Cache
}

// Create new branch
func (u *Branch) Create(ctx context.Context, in *users.Branch) (*users.Branch, error) {
	var output users.Branch
	var err error
	var branchModel model.Branch

	// basic validation
	{
		if len(in.GetName()) == 0 {
			return &output, status.Error(codes.InvalidArgument, "Please supply valid name")
		}

		if len(in.GetAddress()) == 0 {
			return &output, status.Error(codes.InvalidArgument, "Please supply valid company address")
		}

		if len(in.GetCity()) == 0 {
			return &output, status.Error(codes.InvalidArgument, "Please supply valid company city")
		}

		if len(in.GetProvince()) == 0 {
			return &output, status.Error(codes.InvalidArgument, "Please supply valid company province")
		}

		if len(in.GetPhone()) == 0 {
			return &output, status.Error(codes.InvalidArgument, "Please supply valid company phone")
		}

		if len(in.GetPic()) == 0 {
			return &output, status.Error(codes.InvalidArgument, "Please supply valid company pic")
		}

		if len(in.GetPicPhone()) == 0 {
			return &output, status.Error(codes.InvalidArgument, "Please supply valid company pic phone")
		}
	}

	// code validation
	{
		if len(in.GetCode()) == 0 {
			return &output, status.Error(codes.InvalidArgument, "Please supply valid code")
		}

		branchModel = model.Branch{}
		branchModel.Pb.Code = in.GetCode()
		err = branchModel.GetByCode(ctx, u.Db)
		if err != nil {
			if st, ok := status.FromError(err); ok && st.Code() != codes.NotFound {
				return &output, err
			}
		}

		if len(branchModel.Pb.GetId()) > 0 {
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

	if len(userLogin.Pb.GetBranchId()) > 0 {
		return &output, status.Error(codes.PermissionDenied, "Need user region/company to add new branch")
	}

	// region validation
	if len(userLogin.Pb.GetRegionId()) == 0 {
		if len(in.GetRegionId()) == 0 {
			return &output, status.Error(codes.InvalidArgument, "Please supply valid region")
		}

		regionModel := model.Region{}
		regionModel.Pb.Id = in.GetRegionId()
		err = regionModel.Get(ctx, u.Db)
		if err != nil {
			return &output, err
		}

		if regionModel.Pb.GetCompanyId() != userLogin.Pb.GetCompanyId() {
			return &output, status.Error(codes.PermissionDenied, "its not your region")
		}
	} else {
		if len(in.GetRegionId()) > 0 && in.GetRegionId() != userLogin.Pb.GetRegionId() {
			return &output, status.Error(codes.PermissionDenied, "its not your region")
		}

		if (len(in.GetRegionId())) == 0 {
			in.RegionId = userLogin.Pb.GetRegionId()
		}
	}

	branchModel = model.Branch{}
	branchModel.Pb = users.Branch{
		Code:      in.GetCode(),
		CompanyId: in.GetCompanyId(),
		Name:      in.GetName(),
		Address:   in.GetAddress(),
		City:      in.GetCity(),
		Npwp:      in.GetNpwp(),
		Phone:     in.GetPhone(),
		Pic:       in.GetPic(),
		PicPhone:  in.GetPicPhone(),
		Province:  in.GetProvince(),
		RegionId:  in.GetRegionId(),
	}
	tx, err := u.Db.BeginTx(ctx, nil)
	if err != nil {
		return &output, status.Errorf(codes.Internal, "begin tx: %v", err)
	}

	err = branchModel.Create(ctx, u.Db, tx)
	if err != nil {
		tx.Rollback()
		return &output, err
	}

	tx.Commit()

	return &branchModel.Pb, nil
}

// Update branch
func (u *Branch) Update(ctx context.Context, in *users.Branch) (*users.Branch, error) {
	var output users.Branch
	var err error
	var branchModel model.Branch

	// basic validation
	{
		if len(in.GetId()) == 0 {
			return &output, status.Error(codes.InvalidArgument, "Please supply valid id")
		}
		branchModel.Pb.Id = in.GetId()
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

	if len(userLogin.Pb.GetBranchId()) > 0 && userLogin.Pb.GetBranchId() != in.GetId() {
		return &output, status.Error(codes.Unauthenticated, "its not your branch")
	}

	err = branchModel.Get(ctx, u.Db)
	if err != nil {
		return &output, err
	}

	if userLogin.Pb.GetCompanyId() != branchModel.Pb.GetCompanyId() {
		return &output, status.Error(codes.Unauthenticated, "its not your company")
	}

	if len(userLogin.Pb.GetRegionId()) > 0 && userLogin.Pb.GetRegionId() != branchModel.Pb.GetRegionId() {
		return &output, status.Error(codes.Unauthenticated, "its not your region")
	}

	if len(userLogin.Pb.GetRegionId()) == 0 &&
		len(in.GetRegionId()) > 0 &&
		branchModel.Pb.GetRegionId() != in.GetRegionId() {

		regionModel := model.Region{}
		regionModel.Pb.Id = in.GetRegionId()
		err = regionModel.Get(ctx, u.Db)
		if err != nil {
			return &output, err
		}
		branchModel.UpdateRegion = true
		branchModel.Pb.RegionId = in.GetRegionId()
	}

	if len(in.GetName()) > 0 {
		branchModel.Pb.Name = in.GetName()
	}

	if len(in.GetAddress()) == 0 {
		branchModel.Pb.Address = in.GetAddress()
	}

	if len(in.GetCity()) == 0 {
		branchModel.Pb.City = in.GetCity()
	}

	if len(in.GetProvince()) == 0 {
		branchModel.Pb.Province = in.GetProvince()
	}

	if len(in.GetPhone()) == 0 {
		branchModel.Pb.Phone = in.GetPhone()
	}

	if len(in.GetPic()) == 0 {
		branchModel.Pb.Pic = in.GetPic()
	}

	if len(in.GetPicPhone()) == 0 {
		branchModel.Pb.PicPhone = in.GetPicPhone()
	}

	tx, err := u.Db.BeginTx(ctx, nil)
	if err != nil {
		return &output, status.Errorf(codes.Internal, "begin tx: %v", err)
	}

	err = branchModel.Update(ctx, u.Db, tx)
	if err != nil {
		tx.Rollback()
		return &output, err
	}

	tx.Commit()

	return &branchModel.Pb, nil
}

// View branch
func (u *Branch) View(ctx context.Context, in *users.Id) (*users.Branch, error) {
	var output users.Branch
	var err error
	var branchModel model.Branch

	// basic validation
	{
		if len(in.GetId()) == 0 {
			return &output, status.Error(codes.InvalidArgument, "Please supply valid id")
		}
		branchModel.Pb.Id = in.GetId()
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

	err = branchModel.Get(ctx, u.Db)
	if err != nil {
		return &output, err
	}

	if userLogin.Pb.GetCompanyId() != branchModel.Pb.GetCompanyId() {
		return &output, status.Error(codes.Unauthenticated, "its not your company")
	}

	if len(userLogin.Pb.GetRegionId()) > 0 && userLogin.Pb.GetRegionId() != branchModel.Pb.GetRegionId() {
		return &output, status.Error(codes.Unauthenticated, "its not your region")
	}

	if len(userLogin.Pb.GetBranchId()) > 0 && userLogin.Pb.GetBranchId() != branchModel.Pb.GetId() {
		return &output, status.Error(codes.Unauthenticated, "its not your branch")
	}

	return &branchModel.Pb, err
}

// Delete branch
func (u *Branch) Delete(ctx context.Context, in *users.Id) (*users.Boolean, error) {
	var output users.Boolean
	output.Boolean = false
	var err error
	var branchModel model.Branch

	// basic validation
	{
		if len(in.GetId()) == 0 {
			return &output, status.Error(codes.InvalidArgument, "Please supply valid id")
		}
		branchModel.Pb.Id = in.GetId()
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

	err = branchModel.Get(ctx, u.Db)
	if err != nil {
		return &output, err
	}

	if userLogin.Pb.GetCompanyId() != branchModel.Pb.GetCompanyId() {
		return &output, status.Error(codes.Unauthenticated, "its not your company")
	}

	if len(userLogin.Pb.GetRegionId()) > 0 && userLogin.Pb.GetRegionId() != branchModel.Pb.GetRegionId() {
		return &output, status.Error(codes.Unauthenticated, "its not your region")
	}

	err = branchModel.Delete(ctx, u.Db)
	if err != nil {
		return &output, err
	}

	output.Boolean = true
	return &output, err
}

// List branches
func (u *Branch) List(in *users.ListBranchRequest, stream users.BranchService_ListServer) error {
	return nil
}
