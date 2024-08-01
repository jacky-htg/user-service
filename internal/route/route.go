package route

import (
	"database/sql"
	"log"

	"github.com/jacky-htg/erp-pkg/db/redis"
	"github.com/jacky-htg/erp-proto/go/pb/users"
	"github.com/jacky-htg/user-service/internal/service"
	"google.golang.org/grpc"
)

// GrpcRoute func
func GrpcRoute(grpcServer *grpc.Server, db *sql.DB, log *log.Logger, cache *redis.Cache) {
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

	featureServer := service.Feature{Db: db, Cache: cache}
	users.RegisterFeatureServiceServer(grpcServer, &featureServer)

	packageFeatureServer := service.PackageFeature{Db: db, Cache: cache}
	users.RegisterPackageFeatureServiceServer(grpcServer, &packageFeatureServer)

	accessServer := service.Access{Db: db, Cache: cache}
	users.RegisterAccessServiceServer(grpcServer, &accessServer)

	groupServer := service.Group{Db: db, Cache: cache}
	users.RegisterGroupServiceServer(grpcServer, &groupServer)
}
