package service

import (
	"context"
	"database/sql"
	"user-service/internal/model"
	"user-service/internal/pkg/db/redis"
	"user-service/pb/users"
)

// Access struct
type Access struct {
	Db    *sql.DB
	Cache *redis.Cache
	users.UnimplementedAccessServiceServer
}

// List access
func (u *Access) List(ctx context.Context, in *users.MyEmpty) (*users.ListAccessResponse, error) {
	var accessModel model.Access
	return accessModel.List(ctx, u.Db)
}
