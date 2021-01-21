package service

import (
	"context"
	"database/sql"
	"user-service/internal/model"
	"user-service/internal/pkg/app"
	"user-service/internal/pkg/db/redis"
	"user-service/pb/users"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Employee struct
type Employee struct {
	Db    *sql.DB
	Cache *redis.Cache
}

// Create Employee
func (u *Employee) Create(ctx context.Context, in *users.Employee) (*users.Employee, error) {
	var output users.Employee
	var err error
	var employeeModel model.Employee

	// basic validation
	{
		if len(in.GetName()) == 0 {
			return &output, status.Error(codes.InvalidArgument, "Please supply valid name")
		}

		if len(in.GetAddress()) == 0 {
			return &output, status.Error(codes.InvalidArgument, "Please supply valid address")
		}

		if len(in.GetCity()) == 0 {
			return &output, status.Error(codes.InvalidArgument, "Please supply valid city")
		}

		if len(in.GetProvince()) == 0 {
			return &output, status.Error(codes.InvalidArgument, "Please supply valid province")
		}

		if len(in.GetJabatan()) == 0 {
			return &output, status.Error(codes.InvalidArgument, "Please supply valid jabatan")
		}
	}

	// code validation
	{
		if len(in.GetCode()) == 0 {
			return &output, status.Error(codes.InvalidArgument, "Please supply valid code")
		}

		employeeModel = model.Employee{}
		employeeModel.Pb.Code = in.GetCode()
		err = employeeModel.GetByCode(ctx, u.Db)
		if err != nil {
			if st, ok := status.FromError(err); ok && st.Code() != codes.NotFound {
				return &output, err
			}
		}

		if len(employeeModel.Pb.GetId()) > 0 {
			return &output, status.Error(codes.AlreadyExists, "code must be unique")
		}
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

	// user validation
	{
		if len(in.GetUser().GetId()) == 0 {
			return &output, status.Error(codes.InvalidArgument, "Please supply valid user_id")
		}

		userModel := model.User{}
		userModel.Pb.Id = in.GetUser().GetId()
		err = userModel.Get(ctx, u.Db)
		if err != nil {
			return &output, err
		}

		if userLogin.Pb.GetCompanyId() != userModel.Pb.GetCompanyId() {
			return &output, status.Error(codes.PermissionDenied, "its not your company")
		}

		if len(userLogin.Pb.GetRegionId()) > 0 && userLogin.Pb.GetRegionId() != userModel.Pb.GetRegionId() {
			return &output, status.Error(codes.PermissionDenied, "its not your company")
		}

		if len(userLogin.Pb.GetBranchId()) > 0 && userLogin.Pb.GetBranchId() != userModel.Pb.GetBranchId() {
			return &output, status.Error(codes.PermissionDenied, "its not your branch")
		}

		in.User = &userModel.Pb
	}

	employeeModel = model.Employee{}
	employeeModel.Pb = users.Employee{
		Code:     in.GetCode(),
		Name:     in.GetName(),
		Address:  in.GetAddress(),
		City:     in.GetCity(),
		Province: in.GetProvince(),
		Jabatan:  in.GetJabatan(),
		User:     in.GetUser(),
	}
	err = employeeModel.Create(ctx, u.Db)
	if err != nil {
		return &output, err
	}
	return &employeeModel.Pb, nil
}

// Update Employee
func (u *Employee) Update(ctx context.Context, in *users.Employee) (*users.Employee, error) {
	var output users.Employee
	var err error
	var employeeModel model.Employee

	// basic validation
	{
		if len(in.GetId()) == 0 {
			return &output, status.Error(codes.InvalidArgument, "Please supply valid id")
		}
		employeeModel.Pb.Id = in.GetId()
	}

	employeeModel.Pb.Id = in.GetId()
	err = employeeModel.Get(ctx, u.Db)
	if err != nil {
		return &output, err
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

	if len(userLogin.Pb.GetBranchId()) > 0 && userLogin.Pb.GetBranchId() != employeeModel.Pb.GetUser().GetBranchId() {
		return &output, status.Error(codes.Unauthenticated, "its not your branch")
	}

	if userLogin.Pb.GetCompanyId() != employeeModel.Pb.GetUser().GetCompanyId() {
		return &output, status.Error(codes.Unauthenticated, "its not your company")
	}

	if len(userLogin.Pb.GetRegionId()) > 0 && userLogin.Pb.GetRegionId() != employeeModel.Pb.GetUser().GetRegionId() {
		return &output, status.Error(codes.Unauthenticated, "its not your region")
	}

	if len(in.GetUser().GetId()) > 0 && in.GetUser().GetId() != employeeModel.Pb.GetUser().GetId() {
		var userInput model.User
		userInput.Pb.Id = in.GetUser().GetId()
		err = userInput.Get(ctx, u.Db)
		if err != nil {
			return &output, err
		}
		if userLogin.Pb.GetCompanyId() != userInput.Pb.GetCompanyId() {
			return &output, status.Error(codes.PermissionDenied, "its not your company")
		}

		if len(userLogin.Pb.GetRegionId()) > 0 && userLogin.Pb.GetRegionId() != userInput.Pb.GetRegionId() {
			return &output, status.Error(codes.PermissionDenied, "its not your company")
		}

		if len(userLogin.Pb.GetBranchId()) > 0 && userLogin.Pb.GetBranchId() != userInput.Pb.GetBranchId() {
			return &output, status.Error(codes.PermissionDenied, "its not your branch")
		}

		employeeModel.Pb.User = &userInput.Pb
	}

	if len(in.GetName()) > 0 {
		employeeModel.Pb.Name = in.GetName()
	}

	if len(in.GetAddress()) > 0 {
		employeeModel.Pb.Address = in.GetAddress()
	}

	if len(in.GetCity()) > 0 {
		employeeModel.Pb.City = in.GetCity()
	}

	if len(in.GetProvince()) > 0 {
		employeeModel.Pb.Province = in.GetProvince()
	}

	if len(in.GetJabatan()) > 0 {
		employeeModel.Pb.Jabatan = in.GetJabatan()
	}

	err = employeeModel.Update(ctx, u.Db)
	if err != nil {
		return &output, err
	}

	return &employeeModel.Pb, nil
}

// View Employee
func (u *Employee) View(ctx context.Context, in *users.Id) (*users.Employee, error) {
	var output users.Employee
	var err error
	var employeeModel model.Employee

	// basic validation
	{
		if len(in.GetId()) == 0 {
			return &output, status.Error(codes.InvalidArgument, "Please supply valid id")
		}
		employeeModel.Pb.Id = in.GetId()
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

	err = employeeModel.Get(ctx, u.Db)
	if err != nil {
		return &output, err
	}

	if userLogin.Pb.GetCompanyId() != employeeModel.Pb.GetUser().GetCompanyId() {
		return &output, status.Error(codes.Unauthenticated, "its not your company")
	}

	if len(userLogin.Pb.GetRegionId()) > 0 && userLogin.Pb.GetRegionId() != employeeModel.Pb.GetUser().GetRegionId() {
		return &output, status.Error(codes.Unauthenticated, "its not your region")
	}

	if len(userLogin.Pb.GetBranchId()) > 0 && userLogin.Pb.GetBranchId() != employeeModel.Pb.GetUser().GetBranchId() {
		return &output, status.Error(codes.Unauthenticated, "its not your branch")
	}

	return &employeeModel.Pb, nil
}

// Delete Employee
func (u *Employee) Delete(ctx context.Context, in *users.Id) (*users.Boolean, error) {
	var output users.Boolean
	output.Boolean = false
	var err error
	var employeeModel model.Employee

	// basic validation
	{
		if len(in.GetId()) == 0 {
			return &output, status.Error(codes.InvalidArgument, "Please supply valid id")
		}
		employeeModel.Pb.Id = in.GetId()
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

	err = employeeModel.Get(ctx, u.Db)
	if err != nil {
		return &output, err
	}

	if userLogin.Pb.GetCompanyId() != employeeModel.Pb.GetUser().GetCompanyId() {
		return &output, status.Error(codes.Unauthenticated, "its not your company")
	}

	if len(userLogin.Pb.GetRegionId()) > 0 && userLogin.Pb.GetRegionId() != employeeModel.Pb.GetUser().GetRegionId() {
		return &output, status.Error(codes.Unauthenticated, "its not your region")
	}

	if len(userLogin.Pb.GetBranchId()) > 0 && userLogin.Pb.GetBranchId() != employeeModel.Pb.GetUser().GetBranchId() {
		return &output, status.Error(codes.Unauthenticated, "its not your branch")
	}

	err = employeeModel.Delete(ctx, u.Db)
	if err != nil {
		return &output, err
	}

	output.Boolean = true
	return &output, nil
}

// List Employee
func (u *Employee) List(in *users.ListEmployeeRequest, stream users.EmployeeService_ListServer) error {
	ctx := stream.Context()
	ctx, err := getMetadata(ctx)
	if err != nil {
		return err
	}

	// get user login
	var userLogin model.User
	userLogin.Pb.Id = ctx.Value(app.Ctx("userID")).(string)
	err = userLogin.Get(ctx, u.Db)
	if err != nil {
		return err
	}

	if len(in.GetRegionId()) > 0 {
		regionModel := model.Region{}
		regionModel.Pb.Id = in.GetRegionId()
		err = regionModel.Get(ctx, u.Db)
		if err != nil {
			return err
		}

		if regionModel.Pb.GetCompanyId() != ctx.Value(app.Ctx("companyID")).(string) {
			return status.Error(codes.InvalidArgument, "its not your company")
		}
	} else {
		if len(userLogin.Pb.GetRegionId()) > 0 {
			in.RegionId = userLogin.Pb.GetRegionId()
		}
	}

	if len(in.GetBranchId()) > 0 {
		branchModel := model.Branch{}
		branchModel.Pb.Id = in.GetBranchId()
		err = branchModel.Get(ctx, u.Db)
		if err != nil {
			return err
		}

		if branchModel.Pb.GetCompanyId() != ctx.Value(app.Ctx("companyID")).(string) {
			return status.Error(codes.InvalidArgument, "its not your company")
		}

		if branchModel.Pb.GetRegionId() != userLogin.Pb.GetRegionId() {
			return status.Error(codes.InvalidArgument, "its not your region")
		}
	} else {
		if len(userLogin.Pb.GetBranchId()) > 0 {
			in.BranchId = userLogin.Pb.GetBranchId()
		}
	}

	var employeeModel model.Employee
	query, paramQueries, paginationResponse, err := employeeModel.ListQuery(ctx, u.Db, in, &userLogin.Pb)

	rows, err := u.Db.QueryContext(ctx, query, paramQueries...)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}
	defer rows.Close()
	paginationResponse.RegionId = in.GetRegionId()
	paginationResponse.BranchId = in.GetBranchId()
	paginationResponse.Pagination = in.GetPagination()

	for rows.Next() {
		err := contextError(ctx)
		if err != nil {
			return err
		}

		var pbEmployee users.Employee
		var pbUser users.User
		var regionID, branchID sql.NullString
		err = rows.Scan(
			&pbEmployee.Id, &pbEmployee.Name, &pbEmployee.Code, &pbEmployee.Address,
			&pbEmployee.City, &pbEmployee.Province, &pbEmployee.Jabatan,
			&pbUser.Id, &pbUser.CompanyId, &regionID, &branchID, &pbUser.Name, &pbUser.Email,
		)
		if err != nil {
			return err
		}

		pbUser.RegionId = regionID.String
		pbUser.BranchId = branchID.String
		pbEmployee.User = &pbUser

		res := &users.ListEmployeeResponse{
			Pagination: paginationResponse,
			Employee:   &pbEmployee,
		}

		err = stream.Send(res)
		if err != nil {
			return status.Errorf(codes.Unknown, "cannot send stream response: %v", err)
		}
	}

	return nil
}
