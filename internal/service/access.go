package service

import (
	"context"
	"database/sql"

	"github.com/jacky-htg/erp-pkg/db/redis"
	"github.com/jacky-htg/erp-proto/go/pb/users"
	"github.com/jacky-htg/user-service/internal/model"
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
