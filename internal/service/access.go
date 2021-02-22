package service

import (
	"context"
	"database/sql"
	"user-service/internal/model"
	"user-service/internal/pkg/db/redis"
	"user-service/pb/users"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Access struct
type Access struct {
	Db    *sql.DB
	Cache *redis.Cache
}

// List access
func (u *Access) List(ctx context.Context, in *users.MyEmpty) (*users.Access, error) {
	var accessModel model.Access

	tx, err := u.Db.BeginTx(ctx, nil)
	if err != nil {
		return &accessModel.Pb, status.Errorf(codes.Internal, "begin tx: %v", err)
	}

	err = accessModel.GetRoot(ctx, tx, true)
	if err != nil {
		tx.Rollback()
		return &accessModel.Pb, err
	}

	return &accessModel.Pb, nil
}
