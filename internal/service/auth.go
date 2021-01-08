package service

import (
	"context"
	"database/sql"
	"time"
	"user-service/internal/model"

	"github.com/golang/protobuf/ptypes"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"user-service/internal/pkg/app"
	"user-service/internal/pkg/db/redis"
	"user-service/internal/pkg/token"
	users "user-service/pb"
)

// Auth struct
type Auth struct {
	Db    *sql.DB
	Cache *redis.Cache
}

// Login service
func (u *Auth) Login(ctx context.Context, in *users.LoginRequest) (*users.LoginResponse, error) {
	var output users.LoginResponse
	if len(in.GetUsername()) == 0 {
		return &output, status.Error(codes.InvalidArgument, "Please supply valid username")
	}

	if len(in.GetPassword()) == 0 {
		return &output, status.Error(codes.InvalidArgument, "Please supply valid password")
	}

	var userModel model.User
	userModel.Pb.Username = in.GetUsername()
	err := userModel.GetByUserNamePassword(ctx, u.Db, in.GetPassword())
	if err != nil {
		return &output, err
	}

	output.User = &userModel.Pb
	output.Token, err = token.ClaimToken(output.User.GetEmail())
	if err != nil {
		return &output, status.Errorf(codes.Internal, "claim token: %v", err)
	}

	return &output, nil
}

// ForgotPassword service
func (u *Auth) ForgotPassword(ctx context.Context, in *users.ForgotPasswordRequest) (*users.Message, error) {
	var output users.Message
	output.Message = "Failed"
	if len(in.GetEmail()) == 0 {
		return &output, status.Error(codes.InvalidArgument, "Please supply valid email")
	}

	var userModel model.User
	userModel.Pb.Email = in.GetEmail()
	err := userModel.GetByEmail(ctx, u.Db)
	if err != nil {
		return &output, err
	}

	var requestPasswordModel model.RequestPassword
	requestPasswordModel.Pb.UserId = userModel.Pb.GetId()
	err = requestPasswordModel.Create(ctx, u.Db)
	if err != nil {
		return &output, err
	}

	// TODO : send email change password

	output.Message = "Success"
	return &output, nil
}

// ResetPassword service
func (u *Auth) ResetPassword(ctx context.Context, in *users.ResetPasswordRequest) (*users.Message, error) {
	output := users.Message{Message: "Failed"}

	if len(in.GetToken()) == 0 {
		return &output, status.Error(codes.InvalidArgument, "Please supply valid token")
	}

	if len(in.GetNewPassword()) == 0 {
		return &output, status.Error(codes.InvalidArgument, "Please supply valid new password")
	}

	if len(in.GetRePassword()) == 0 {
		return &output, status.Error(codes.InvalidArgument, "Please supply valid re password")
	}

	if in.GetNewPassword() != in.GetRePassword() {
		return &output, status.Error(codes.InvalidArgument, "new password not match with re password")
	}

	var requestPasswordModel model.RequestPassword
	requestPasswordModel.Pb.Id = in.GetToken()
	err := requestPasswordModel.Get(ctx, u.Db)
	if err != nil {
		return &output, err
	}

	if requestPasswordModel.Pb.GetIsUsed() {
		return &output, status.Error(codes.PermissionDenied, "token has been used")
	}

	createdAt, err := ptypes.Timestamp(requestPasswordModel.Pb.GetCreatedAt())
	if err != nil {
		return &output, status.Errorf(codes.Internal, "ptypes timestamp: %v", err)
	}

	if time.Now().UTC().After(createdAt.Add(time.Hour * 2 * 24)) {
		return &output, status.Error(codes.PermissionDenied, "token has been expired")
	}

	pass, err := bcrypt.GenerateFromPassword([]byte(in.GetNewPassword()), bcrypt.DefaultCost)
	if err != nil {
		return &output, status.Errorf(codes.Internal, "hash password: %v", err)
	}

	var userModel model.User
	userModel.Pb.Id = requestPasswordModel.Pb.GetUserId()
	tx, err := u.Db.BeginTx(ctx, nil)
	if err != nil {
		return &output, status.Errorf(codes.Internal, "begin tx: %v", err)
	}

	err = userModel.ChangePassword(ctx, tx, string(pass))
	if err != nil {
		tx.Rollback()
		return &output, err
	}

	err = requestPasswordModel.UpdateIsUsed(ctx, tx)
	if err != nil {
		tx.Rollback()
		return &output, err
	}

	tx.Commit()

	output.Message = "Success"

	return &output, nil
}

// ChangePassword service
func (u *Auth) ChangePassword(ctx context.Context, in *users.ChangePasswordRequest) (*users.Message, error) {
	var output users.Message
	output.Message = "Failed"

	ctx, err := getMetadata(ctx)
	if err != nil {
		return &output, err
	}

	if len(in.GetOldPassword()) == 0 {
		return &output, status.Error(codes.InvalidArgument, "Please supply valid current password")
	}

	if len(in.GetNewPassword()) == 0 {
		return &output, status.Error(codes.InvalidArgument, "Please supply valid new password")
	}

	if len(in.GetRePassword()) == 0 {
		return &output, status.Error(codes.InvalidArgument, "Please supply valid re password")
	}

	if in.GetNewPassword() != in.GetRePassword() {
		return &output, status.Error(codes.InvalidArgument, "new password not match with re password")
	}

	var userModel model.User
	userModel.Pb.Id = ctx.Value(app.Ctx("userID")).(string)
	err = userModel.GetByPassword(ctx, u.Db, in.GetOldPassword())
	if err != nil {
		return &output, err
	}

	pass, err := bcrypt.GenerateFromPassword([]byte(in.GetNewPassword()), bcrypt.DefaultCost)
	if err != nil {
		return &output, status.Errorf(codes.Internal, "hash password: %v", err)
	}

	tx, err := u.Db.BeginTx(ctx, nil)
	if err != nil {
		return &output, status.Errorf(codes.Internal, "begin tx: %v", err)
	}

	err = userModel.ChangePassword(ctx, tx, string(pass))
	if err != nil {
		tx.Rollback()
		return &output, err
	}

	tx.Commit()
	output.Message = "success"

	return &output, nil
}

// IsAuth service
func (u *Auth) IsAuth(ctx context.Context, in *users.String) (*users.Boolean, error) {
	output := users.Boolean{Boolean: false}

	ctx, err := getMetadata(ctx)
	if err != nil {
		return &output, err
	}

	if len(in.GetString_()) == 0 {
		return &output, status.Error(codes.InvalidArgument, "Please supply valid access")
	}

	var userModel model.User
	userModel.Pb.Id = ctx.Value(app.Ctx("userID")).(string)
	err = userModel.IsAuth(ctx, u.Db, in.GetString_())
	if err != nil {
		return &output, err
	}

	output.Boolean = true
	return &output, nil
}

func getMetadataToken(ctx context.Context) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	token := md["token"]
	if len(token) == 0 {
		return ctx, status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}

	ctx = context.WithValue(ctx, app.Ctx("token"), token[0])

	return ctx, nil
}

func getMetadata(ctx context.Context) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	userID := md["user_id"]
	if len(userID) == 0 {
		return ctx, status.Errorf(codes.Unauthenticated, "user_id is not provided")
	}

	ctx = context.WithValue(ctx, app.Ctx("userID"), userID[0])

	companyID := md["company_id"]
	if len(companyID) == 0 {
		return ctx, status.Errorf(codes.Unauthenticated, "company_id is not provided")
	}

	ctx = context.WithValue(ctx, app.Ctx("companyID"), companyID[0])

	return ctx, nil
}
