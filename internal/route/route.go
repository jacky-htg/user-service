package route

import (
	"database/sql"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"user-service/internal/pkg/db/redis"
	"user-service/internal/service"
	users "user-service/pb"
)

// GrpcRoute func
func GrpcRoute(grpcServer *grpc.Server, db *sql.DB, log *logrus.Entry, cache *redis.Cache) {
	authServer := service.Auth{Db: db, Cache: cache}
	users.RegisterAuthServiceServer(grpcServer, &authServer)

	userServer := service.User{Db: db, Cache: cache}
	users.RegisterUserServiceServer(grpcServer, &userServer)

	companyServer := service.Company{Db: db, Cache: cache}
	users.RegisterCompanyServiceServer(grpcServer, &companyServer)

	regionServer := service.Region{Db: db, Cache: cache}
	users.RegisterRegionServiceServer(grpcServer, &regionServer)

	branchServer := service.Branch{Db: db, Cache: cache}
	users.RegisterBranchServiceServer(grpcServer, &branchServer)

	employeeServer := service.Employee{Db: db, Cache: cache}
	users.RegisterEmployeeServiceServer(grpcServer, &employeeServer)
}
