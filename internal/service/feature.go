package service

import (
	"database/sql"
	"user-service/internal/pkg/db/redis"
	"user-service/pb/users"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Feature struct
type Feature struct {
	Db    *sql.DB
	Cache *redis.Cache
	users.UnimplementedFeatureServiceServer
}

// List feature
func (u *Feature) List(in *users.MyEmpty, stream users.FeatureService_ListServer) error {
	ctx := stream.Context()
	rows, err := u.Db.QueryContext(ctx, `SELECT id, name from features`)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		err := contextError(ctx)
		if err != nil {
			return err
		}

		var pbFeature users.Feature
		err = rows.Scan(&pbFeature.Id, &pbFeature.Name)
		if err != nil {
			return err
		}

		err = stream.Send(&users.ListFeatureResponse{Feature: &pbFeature})
		if err != nil {
			return status.Errorf(codes.Unknown, "cannot send stream response: %v", err)
		}
	}
	return nil
}
