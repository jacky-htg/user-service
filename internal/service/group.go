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
			return groupModel.Pb, status.Error(codes.InvalidArgument, "Please supply valid name")
		}
	}

	ctx, err = getMetadata(ctx)
	if err != nil {
		return groupModel.Pb, err
	}

	// company validation
	{
		if len(in.GetCompanyId()) > 0 && in.GetCompanyId() != ctx.Value(app.Ctx("companyID")).(string) {
			return groupModel.Pb, status.Error(codes.PermissionDenied, "Please supply valid company id")
		}
		in.CompanyId = ctx.Value(app.Ctx("companyID")).(string)
	}

	groupModel.Pb = in
	err = groupModel.Create(ctx, u.Db)
	if err != nil {
		return groupModel.Pb, err
	}

	return groupModel.Pb, nil
}

// Update Group
func (u *Group) Update(ctx context.Context, in *users.Group) (*users.Group, error) {
	var groupModel model.Group

	return groupModel.Pb, nil
}

// View Group
func (u *Group) View(ctx context.Context, in *users.Id) (*users.Group, error) {
	var groupModel model.Group

	return groupModel.Pb, nil
}

// Delete Group
func (u *Group) Delete(ctx context.Context, in *users.Id) (*users.Boolean, error) {
	var output users.Boolean
	output.Boolean = false

	output.Boolean = true
	return &output, nil
}

// List Group
func (u *Group) List(in *users.ListGroupRequest, stream users.GroupService_ListServer) error {
	return nil
}

// GrantAccess Group
func (u *Group) GrantAccess(ctx context.Context, in *users.GrantAccessRequest) (*users.Message, error) {
	var output users.Message
	output.Message = "Failed"

	output.Message = "Success"
	return &output, nil
}

// RevokeAccess Group
func (u *Group) RevokeAccess(ctx context.Context, in *users.GrantAccessRequest) (*users.Message, error) {
	var output users.Message
	output.Message = "Failed"

	output.Message = "Success"
	return &output, nil
}
