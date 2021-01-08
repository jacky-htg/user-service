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
}
