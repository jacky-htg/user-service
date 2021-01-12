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
	return &users.Employee{}, nil
}

// View Employee
func (u *Employee) View(ctx context.Context, in *users.Id) (*users.Employee, error) {
	return &users.Employee{}, nil
}

// Delete Employee
func (u *Employee) Delete(ctx context.Context, in *users.Id) (*users.Boolean, error) {
	return &users.Boolean{}, nil
}

// List Employee
func (u *Employee) List(in *users.ListEmployeeRequest, stream users.EmployeeService_ListServer) error {
	return nil
}
