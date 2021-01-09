package service

import (
	"context"
	"database/sql"
	"regexp"
	"user-service/internal/model"
	"user-service/internal/pkg/app"
	"user-service/internal/pkg/db/redis"
	users "user-service/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// User struct
type User struct {
	Db    *sql.DB
	Cache *redis.Cache
}

// Create func
func (u *User) Create(ctx context.Context, in *users.User) (*users.User, error) {
	var output users.User
	var err error
	var userModel model.User

	if in == nil {
		return &output, status.Error(codes.InvalidArgument, "Please supply valid argument")
	}

	if len(in.GetEmail()) == 0 {
		return &output, status.Error(codes.InvalidArgument, "Please supply valid email")
	}

	if len(in.GetGroup().GetId()) == 0 {
		return &output, status.Error(codes.InvalidArgument, "Please supply valid group_id")
	}

	if len(in.GetName()) == 0 {
		return &output, status.Error(codes.InvalidArgument, "Please supply valid name")
	}

	// username validation
	{
		if len(in.GetUsername()) == 0 {
			return &output, status.Error(codes.InvalidArgument, "Please supply valid username")
		}

		if len(in.GetUsername()) < 8 {
			return &output, status.Error(codes.InvalidArgument, "username min 8 character")
		}

		userModel = model.User{}
		userModel.Pb.Username = in.GetUsername()
		err = userModel.GetByUsername(ctx, u.Db)
		if err != nil {
			if st, ok := status.FromError(err); ok && st.Code() != codes.NotFound {
				return &output, err
			}
		}

		if len(userModel.Pb.GetId()) > 0 {
			return &output, status.Error(codes.AlreadyExists, "username must be unique")
		}
	}

	// email validation
	{
		var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
		if valid := func(e string) bool {
			if len(e) < 3 && len(e) > 254 {
				return false
			}
			return emailRegex.MatchString(e)
		}(in.GetEmail()); !valid {
			return &output, status.Error(codes.InvalidArgument, "Please supply valid email")
		}

		userModel = model.User{}
		userModel.Pb.Email = in.GetEmail()
		err = userModel.GetByEmail(ctx, u.Db)
		if err != nil {
			if st, ok := status.FromError(err); ok && st.Code() != codes.NotFound {
				return &output, err
			}
		}

		if len(userModel.Pb.GetId()) > 0 {
			return &output, status.Error(codes.AlreadyExists, "email must be unique")
		}
	}

	ctx, err = getMetadata(ctx)
	if err != nil {
		return &output, err
	}

	// company validation
	if len(in.GetCompanyId()) > 0 && in.GetCompanyId() != ctx.Value(app.Ctx("companyID")).(string) {
		return &output, status.Error(codes.PermissionDenied, "Please supply valid company id")
	}
	in.CompanyId = ctx.Value(app.Ctx("companyID")).(string)

	var userLogin model.User
	userLogin.Pb.Id = ctx.Value(app.Ctx("userID")).(string)
	err = userLogin.Get(ctx, u.Db)
	if err != nil {
		return &output, err
	}

	// region validation
	var regionModel model.Region
	if len(in.GetRegionId()) > 0 {
		if len(userLogin.Pb.GetRegionId()) == 0 {
			// check is region belongsto company
			regionModel.Pb.Id = in.GetRegionId()
			err = regionModel.Get(ctx, u.Db)
			if err != nil {
				return &output, err
			}

			if regionModel.Pb.GetCompanyId() != in.GetCompanyId() {
				return &output, status.Error(codes.PermissionDenied, "Please supply valid region id")
			}
		} else {
			if in.GetRegionId() != userLogin.Pb.GetRegionId() {
				return &output, status.Error(codes.PermissionDenied, "Please supply valid region id")
			}
		}
	} else {
		if len(userLogin.Pb.GetRegionId()) > 0 {
			in.RegionId = userLogin.Pb.GetRegionId()
		}
	}

	// branch validation
	if len(in.GetBranchId()) > 0 {
		if len(userLogin.Pb.GetBranchId()) == 0 {
			// check is branch belongsto region
			var branchModel model.Branch
			branchModel.Pb.Id = in.GetBranchId()
			err = branchModel.Get(ctx, u.Db)
			if err != nil {
				return &output, err
			}

			if branchModel.Pb.GetRegionId() != in.GetRegionId() {
				return &output, status.Error(codes.PermissionDenied, "Please supply valid branch id")
			}
		} else {
			if in.GetBranchId() != userLogin.Pb.GetBranchId() {
				return &output, status.Error(codes.PermissionDenied, "Please supply valid branch id")
			}
		}
	} else {
		if len(userLogin.Pb.GetBranchId()) > 0 {
			in.BranchId = userLogin.Pb.GetBranchId()
		}
	}

	// group validation
	var groupModel model.Group
	groupModel.Pb.Id = in.GetGroup().GetId()
	err = groupModel.Get(ctx, u.Db)
	if err != nil {
		return &output, err
	}

	if groupModel.Pb.GetCompanyId() != in.GetCompanyId() {
		return &output, status.Error(codes.PermissionDenied, "Please supply valid group id")
	}

	userModel = model.User{}
	userModel.Pb = users.User{
		BranchId:  in.GetBranchId(),
		CompanyId: in.GetCompanyId(),
		Email:     in.GetEmail(),
		Group:     in.GetGroup(),
		Name:      in.GetName(),
		RegionId:  in.GetRegionId(),
		Username:  in.GetUsername(),
	}
	// TODO generate random password
	userModel.Password = "1234"
	err = userModel.Create(ctx, u.Db)
	if err != nil {
		return &output, err
	}

	// TODO : send email to inform username and password

	return &userModel.Pb, nil
}

// Update func
func (u *User) Update(ctx context.Context, in *users.User) (*users.User, error) {
	var output users.User

	return &output, nil
}

// View func
func (u *User) View(ctx context.Context, in *users.Id) (*users.User, error) {
	var output users.User

	return &output, nil
}

// Delete func
func (u *User) Delete(ctx context.Context, in *users.Id) (*users.Boolean, error) {
	var output users.Boolean
	output.Boolean = false

	output.Boolean = true
	return &output, nil
}

// List func
func (u *User) List(ctx context.Context, in *users.ListUserRequest) (*users.ListUserResponse, error) {
	var output users.ListUserResponse

	return &output, nil
}

// GetByToken func
func (u *User) GetByToken(ctx context.Context, in *users.Empty) (*users.User, error) {
	var output users.User

	return &output, nil
}
