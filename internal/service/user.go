package service

import (
	"context"
	"database/sql"
	"regexp"
	"user-service/internal/model"
	"user-service/internal/pkg/app"
	"user-service/internal/pkg/db/redis"
	"user-service/internal/pkg/token"
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

	// basic validation
	{
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

	// region validation
	{
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
	}

	// branch validation
	{
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
	}

	// group validation
	{
		var groupModel model.Group
		groupModel.Pb.Id = in.GetGroup().GetId()
		err = groupModel.Get(ctx, u.Db)
		if err != nil {
			return &output, err
		}

		if groupModel.Pb.GetCompanyId() != in.GetCompanyId() {
			return &output, status.Error(codes.PermissionDenied, "Please supply valid group id")
		}
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
	var err error
	var userModel model.User

	// basic validation
	{
		if len(in.GetId()) == 0 {
			return &output, status.Error(codes.InvalidArgument, "Please supply valid id")
		}
		userModel.Pb.Id = in.GetId()
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

	// region validation
	if len(in.GetRegionId()) > 0 {
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
	}

	// branch validation
	if len(in.GetBranchId()) > 0 {
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
	}

	// group validation
	if len(in.GetGroup().GetId()) > 0 {
		var groupModel model.Group
		groupModel.Pb.Id = in.GetGroup().GetId()
		err = groupModel.Get(ctx, u.Db)
		if err != nil {
			return &output, err
		}

		if groupModel.Pb.GetCompanyId() != in.GetCompanyId() {
			return &output, status.Error(codes.PermissionDenied, "Please supply valid group id")
		}
	}

	err = userModel.Get(ctx, u.Db)
	if err != nil {
		return &output, err
	}

	err = u.checkFilteringContent(ctx, &userLogin, &userModel)
	if err != nil {
		return &output, err
	}

	if len(in.GetName()) > 0 {
		userModel.Pb.Name = in.GetName()
	}

	if len(in.GetRegionId()) > 0 {
		userModel.Pb.RegionId = in.GetRegionId()
	}

	if len(in.GetBranchId()) > 0 {
		userModel.Pb.BranchId = in.GetBranchId()
	}

	if len(in.GetGroup().GetId()) > 0 {
		userModel.Pb.Name = in.GetGroup().GetId()
	}

	err = userModel.Update(ctx, u.Db)
	if err != nil {
		return &output, err
	}

	return &userModel.Pb, nil
}

// View func
func (u *User) View(ctx context.Context, in *users.Id) (*users.User, error) {
	var output users.User
	var err error
	var userModel model.User

	if len(in.GetId()) == 0 {
		return &output, status.Error(codes.InvalidArgument, "Please supply valid id")
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

	userModel.Pb.Id = in.GetId()
	err = userModel.Get(ctx, u.Db)
	if err != nil {
		return &output, err
	}

	err = u.checkFilteringContent(ctx, &userLogin, &userModel)
	if err != nil {
		return &output, err
	}

	return &userModel.Pb, nil
}

// Delete func
func (u *User) Delete(ctx context.Context, in *users.Id) (*users.Boolean, error) {
	var output users.Boolean
	output.Boolean = false
	var err error
	var userModel model.User

	if len(in.GetId()) == 0 {
		return &output, status.Error(codes.InvalidArgument, "Please supply valid id")
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

	userModel.Pb.Id = in.GetId()
	err = userModel.Get(ctx, u.Db)
	if err != nil {
		return &output, err
	}

	err = u.checkFilteringContent(ctx, &userLogin, &userModel)
	if err != nil {
		return &output, err
	}

	err = userModel.Delete(ctx, u.Db)
	if err != nil {
		return &output, err
	}

	output.Boolean = true
	return &output, nil
}

// List func
func (u *User) List(in *users.ListUserRequest, stream users.UserService_ListServer) error {
	ctx := stream.Context()
	var userModel model.User
	query, paramQueries, paginationResponse, err := userModel.ListQuery(ctx, u.Db, in)

	rows, err := u.Db.QueryContext(ctx, query, paramQueries...)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}
	defer rows.Close()
	paginationResponse.BranchId = in.GetBranchId()
	paginationResponse.CompanyId = in.GetCompanyId()
	paginationResponse.Pagination = in.GetPagination()

	for rows.Next() {
		err := contextError(ctx)
		if err != nil {
			return err
		}

		var pbUser users.User
		var pbGroup users.Group
		var regionID, branchID sql.NullString
		err = rows.Scan(
			&pbUser.Id, &pbUser.CompanyId, &regionID, &branchID, &pbUser.Name, &pbUser.Email,
			&pbGroup.Id, &pbGroup.Name,
		)
		if err != nil {
			return err
		}

		pbUser.RegionId = regionID.String
		pbUser.BranchId = branchID.String
		pbUser.Group = &pbGroup

		res := &users.ListUserResponse{
			Pagination: paginationResponse,
			User:       &pbUser,
		}

		err = stream.Send(res)
		if err != nil {
			return status.Errorf(codes.Unknown, "cannot send stream response: %v", err)
		}
	}

	return nil
}

// GetByToken func
func (u *User) GetByToken(ctx context.Context, in *users.Empty) (*users.User, error) {
	var output users.User
	var err error
	var userModel model.User

	ctx, err = getMetadataToken(ctx)
	if err != nil {
		return &output, err
	}

	// validate token
	isValid, email := token.ValidateToken(ctx.Value(app.Ctx("token")).(string))
	if !isValid {
		return &output, status.Error(codes.Unauthenticated, "invalid token")
	}
	userModel.Pb.Email = email
	err = userModel.GetByEmail(ctx, u.Db)
	if err != nil {
		return &output, err
	}

	return &userModel.Pb, nil
}

func (u *User) checkFilteringContent(ctx context.Context, userLogin *model.User, userModel *model.User) error {
	var err error
	if userModel.Pb.GetCompanyId() != userLogin.Pb.GetCompanyId() {
		return status.Error(codes.PermissionDenied, "user not in your company")
	}

	if len(userLogin.Pb.GetRegionId()) > 0 && userModel.Pb.GetRegionId() != userLogin.Pb.GetRegionId() {
		return status.Error(codes.PermissionDenied, "user not in your region")
	}

	if len(userLogin.Pb.GetRegionId()) == 0 && len(userModel.Pb.GetRegionId()) > 0 {
		// check is region belongsto company
		var regionModel model.Region
		regionModel.Pb.Id = userModel.Pb.GetRegionId()
		err = regionModel.Get(ctx, u.Db)
		if err != nil {
			return err
		}
		if regionModel.Pb.GetCompanyId() != userLogin.Pb.GetCompanyId() {
			return status.Error(codes.PermissionDenied, "user not in your region")
		}
	}

	if len(userLogin.Pb.GetBranchId()) > 0 && userModel.Pb.GetBranchId() != userLogin.Pb.GetBranchId() {
		return status.Error(codes.PermissionDenied, "user not in your branch")
	}

	if len(userLogin.Pb.GetBranchId()) == 0 && len(userModel.Pb.GetBranchId()) > 0 {
		var branchModel model.Branch
		branchModel.Pb.Id = userModel.Pb.GetBranchId()
		err = branchModel.Get(ctx, u.Db)
		if err != nil {
			return err
		}

		if len(userLogin.Pb.GetRegionId()) > 0 && branchModel.Pb.GetRegionId() != userLogin.Pb.GetRegionId() {
			return status.Error(codes.PermissionDenied, "user not in your branch")
		}

		if branchModel.Pb.GetCompanyId() != userLogin.Pb.GetCompanyId() {
			return status.Error(codes.PermissionDenied, "user not in your branch")
		}
	}

	return nil
}
