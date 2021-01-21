package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"user-service/internal/model"
	"user-service/internal/pkg/app"
	"user-service/internal/pkg/db/redis"
	"user-service/pb/users"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Group struct
type Group struct {
	Db    *sql.DB
	Cache *redis.Cache
}

// Create Group
func (u *Group) Create(ctx context.Context, in *users.Group) (*users.Group, error) {
	var groupModel model.Group
	var err error

	// basic validation
	{
		if len(in.GetName()) == 0 {
			return &groupModel.Pb, status.Error(codes.InvalidArgument, "Please supply valid name")
		}
	}

	ctx, err = getMetadata(ctx)
	if err != nil {
		return &groupModel.Pb, err
	}

	// company validation
	{
		if len(in.GetCompanyId()) > 0 && in.GetCompanyId() != ctx.Value(app.Ctx("companyID")).(string) {
			return &groupModel.Pb, status.Error(codes.PermissionDenied, "Please supply valid company id")
		}
		in.CompanyId = ctx.Value(app.Ctx("companyID")).(string)
	}

	groupModel.Pb = users.Group{
		CompanyId: in.GetCompanyId(),
		Name:      in.GetName(),
	}
	err = groupModel.Create(ctx, u.Db)
	if err != nil {
		return &groupModel.Pb, err
	}

	return &groupModel.Pb, nil
}

// Update Group
func (u *Group) Update(ctx context.Context, in *users.Group) (*users.Group, error) {
	var groupModel model.Group
	var err error

	// basic validation
	{
		if len(in.GetId()) == 0 {
			return &groupModel.Pb, status.Error(codes.InvalidArgument, "Please supply valid id")
		}
		groupModel.Pb.Id = in.GetId()
	}

	ctx, err = getMetadata(ctx)
	if err != nil {
		return &groupModel.Pb, err
	}

	// get user login
	var userLogin model.User
	userLogin.Pb.Id = ctx.Value(app.Ctx("userID")).(string)
	err = userLogin.Get(ctx, u.Db)
	if err != nil {
		return &groupModel.Pb, err
	}

	if len(in.GetCompanyId()) > 0 && userLogin.Pb.GetCompanyId() != in.GetCompanyId() {
		return &groupModel.Pb, status.Error(codes.Unauthenticated, "its not your company")
	}

	err = groupModel.Get(ctx, u.Db)
	if err != nil {
		return &groupModel.Pb, err
	}

	if len(in.GetName()) > 0 {
		groupModel.Pb.Name = in.GetName()
	}

	err = groupModel.Update(ctx, u.Db)
	if err != nil {
		return &groupModel.Pb, err
	}

	return &groupModel.Pb, nil
}

// View Group
func (u *Group) View(ctx context.Context, in *users.Id) (*users.Group, error) {
	var groupModel model.Group
	var err error

	// basic validation
	{
		if len(in.GetId()) == 0 {
			return &groupModel.Pb, status.Error(codes.InvalidArgument, "Please supply valid id")
		}
		groupModel.Pb.Id = in.GetId()
	}

	ctx, err = getMetadata(ctx)
	if err != nil {
		return &groupModel.Pb, err
	}

	// get user login
	var userLogin model.User
	userLogin.Pb.Id = ctx.Value(app.Ctx("userID")).(string)
	err = userLogin.Get(ctx, u.Db)
	if err != nil {
		return &groupModel.Pb, err
	}

	err = groupModel.Get(ctx, u.Db)
	if err != nil {
		return &groupModel.Pb, err
	}

	if userLogin.Pb.GetCompanyId() != groupModel.Pb.GetCompanyId() {
		return &groupModel.Pb, status.Error(codes.Unauthenticated, "its not your company")
	}

	return &groupModel.Pb, nil
}

// Delete Group
func (u *Group) Delete(ctx context.Context, in *users.Id) (*users.Boolean, error) {
	var output users.Boolean
	output.Boolean = false

	var groupModel model.Group
	var err error

	// basic validation
	{
		if len(in.GetId()) == 0 {
			return &output, status.Error(codes.InvalidArgument, "Please supply valid id")
		}
		groupModel.Pb.Id = in.GetId()
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

	err = groupModel.Get(ctx, u.Db)
	if err != nil {
		return &output, err
	}

	if userLogin.Pb.GetCompanyId() != groupModel.Pb.GetCompanyId() {
		return &output, status.Error(codes.Unauthenticated, "its not your company")
	}

	err = groupModel.Delete(ctx, u.Db)
	if err != nil {
		return &output, err
	}

	output.Boolean = true
	return &output, nil
}

// List Group
func (u *Group) List(in *users.ListGroupRequest, stream users.GroupService_ListServer) error {
	ctx := stream.Context()
	ctx, err := getMetadata(ctx)
	if err != nil {
		return err
	}

	var groupModel model.Group
	query, paramQueries, paginationResponse, err := groupModel.ListQuery(ctx, u.Db, in)

	rows, err := u.Db.QueryContext(ctx, query, paramQueries...)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}
	defer rows.Close()
	paginationResponse.Pagination = in.GetPagination()

	for rows.Next() {
		err := contextError(ctx)
		if err != nil {
			return err
		}

		var tmpAccess string
		var pbGroup users.Group
		err = rows.Scan(&pbGroup.Id, &pbGroup.CompanyId, &pbGroup.Name, &pbGroup.IsMutable, &tmpAccess)
		if err != nil {
			return err
		}

		err = json.Unmarshal([]byte(tmpAccess), &pbGroup.Access)
		if err != nil {
			return status.Errorf(codes.Internal, "unmarshal access: %v", err)
		}

		res := &users.ListGroupResponse{
			Pagination: paginationResponse,
			Group:      &pbGroup,
		}

		err = stream.Send(res)
		if err != nil {
			return status.Errorf(codes.Unknown, "cannot send stream response: %v", err)
		}
	}
	return nil
}

// GrantAccess Group
func (u *Group) GrantAccess(ctx context.Context, in *users.GrantAccessRequest) (*users.Message, error) {
	var output users.Message
	output.Message = "Failed"

	var err error

	// basic validation
	{
		if len(in.GetGroupId()) == 0 {
			return &output, status.Error(codes.InvalidArgument, "Please supply valid group")
		}

		if len(in.GetAccessId()) == 0 {
			return &output, status.Error(codes.InvalidArgument, "Please supply valid access")
		}
	}

	ctx, err = getMetadata(ctx)
	if err != nil {
		return &output, err
	}

	// group validation
	{
		groupModel := model.Group{}
		groupModel.Pb.Id = in.GetGroupId()
		err = groupModel.Get(ctx, u.Db)
		if err != nil {
			return &output, err
		}

		if groupModel.Pb.GetCompanyId() != ctx.Value(app.Ctx("companyID")).(string) {
			return &output, status.Error(codes.PermissionDenied, "its not your company")
		}
	}

	// access validation
	{
		accessModel := model.Access{}
		accessModel.Pb.Id = in.GetAccessId()
		err = accessModel.Get(ctx, u.Db)
		if err != nil {
			return &output, err
		}
	}

	accessGroup := model.AccessGroup{}
	accessGroup.AccessID = in.GetAccessId()
	accessGroup.GroupID = in.GetGroupId()
	err = accessGroup.Get(ctx, u.Db)
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() != codes.NotFound {
			return &output, err
		}
	} else {
		return &output, status.Error(codes.AlreadyExists, "Already grant access")
	}

	err = accessGroup.Grant(ctx, u.Db)
	if err != nil {
		return &output, err
	}

	output.Message = "Success"
	return &output, nil
}

// RevokeAccess Group
func (u *Group) RevokeAccess(ctx context.Context, in *users.GrantAccessRequest) (*users.Message, error) {
	var output users.Message
	output.Message = "Failed"

	var err error

	// basic validation
	{
		if len(in.GetGroupId()) == 0 {
			return &output, status.Error(codes.InvalidArgument, "Please supply valid group")
		}

		if len(in.GetAccessId()) == 0 {
			return &output, status.Error(codes.InvalidArgument, "Please supply valid access")
		}
	}

	ctx, err = getMetadata(ctx)
	if err != nil {
		return &output, err
	}

	// group validation
	{
		groupModel := model.Group{}
		groupModel.Pb.Id = in.GetGroupId()
		err = groupModel.Get(ctx, u.Db)
		if err != nil {
			return &output, err
		}

		if groupModel.Pb.GetCompanyId() != ctx.Value(app.Ctx("companyID")).(string) {
			return &output, status.Error(codes.PermissionDenied, "its not your company")
		}
	}

	// access validation
	{
		accessModel := model.Access{}
		accessModel.Pb.Id = in.GetAccessId()
		err = accessModel.Get(ctx, u.Db)
		if err != nil {
			return &output, err
		}
	}

	accessGroup := model.AccessGroup{}
	accessGroup.AccessID = in.GetAccessId()
	accessGroup.GroupID = in.GetGroupId()
	err = accessGroup.Get(ctx, u.Db)
	if err != nil {
		return &output, err
	}

	err = accessGroup.Revoke(ctx, u.Db)
	if err != nil {
		return &output, err
	}

	output.Message = "Success"
	return &output, nil
}
