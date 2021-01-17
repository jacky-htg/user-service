package main

import (
	"log"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"user-service/client/service"
	users "user-service/pb"
)

func main() {

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":8000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	ctx := context.Background()

	auth := users.NewAuthServiceClient(conn)
	// user := users.NewUserServiceClient(conn)
	// company := users.NewCompanyServiceClient(conn)
	// region := users.NewRegionServiceClient(conn)
	branch := users.NewBranchServiceClient(conn)
	// employee := users.NewEmployeeServiceClient(conn)
	// feature := users.NewFeatureServiceClient(conn)
	// packageFeature := users.NewPackageFeatureServiceClient(conn)
	// access := users.NewAccessServiceClient(conn)
	// group := users.NewGroupServiceClient(conn)

	ctx = service.Login(ctx, auth)
	// service.ForgotPassword(ctx, auth)
	// service.ResetPassword(ctx, auth)
	// service.ChangePassword(ctx, auth)
	// service.IsAuth(ctx, auth)
	// service.CreateUser(ctx, user)
	// service.UpdateUser(ctx, user)
	// service.ViewUser(ctx, user)
	// service.DeleteUser(ctx, user)
	// service.GetUserByToken(ctx, user)
	// service.ListUser(ctx, user)
	// service.Registration(ctx, company)
	// service.UpdateCompany(ctx, company)
	// service.ViewCompany(ctx, company)
	// service.CreateRegion(ctx, region)
	// service.UpdateRegion(ctx, region)
	// service.ViewRegion(ctx, region)
	// service.DeleteRegion(ctx, region)
	// service.ListRegion(ctx, region)
	service.CreateBranch(ctx, branch)
	// service.UpdateBranch(ctx, branch)
	// service.ViewBranch(ctx, branch)
	// service.DeleteBranch(ctx, branch)
	// service.ListBranch(ctx, branch)
	// service.CreateEmployee(ctx, employee)
	// service.UpdateEmployee(ctx, employee)
	// service.ViewEmployee(ctx, employee)
	// service.DeleteEmployee(ctx, employee)
	// service.ListEmployee(ctx, employee)
	// service.ListFeature(ctx, feature)
	// service.ListPackageFeature(ctx, packageFeature)
	// service.ViewPackageFeature(ctx, packageFeature)
	// service.ViewAccessTree(ctx, access)
	// service.CreateGroup(ctx, group)
	// service.ViewGroup(ctx, group)
	// service.DeleteGroup(ctx, group)
	// service.ListGroup(ctx, group)
	// service.GrantAccess(ctx, group)
	// service.RevokeAccess(ctx, group)
}
