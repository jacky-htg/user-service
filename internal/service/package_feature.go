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

// PackageFeature struct
type PackageFeature struct {
	Db    *sql.DB
	Cache *redis.Cache
}

// View Package Feature
func (u *PackageFeature) View(ctx context.Context, in *users.Id) (*users.PackageOfFeature, error) {
	var output users.PackageOfFeature
	var packageFeatureModel model.FeaturePackage

	// basic validation
	{
		if len(in.GetId()) == 0 {
			return &output, status.Error(codes.InvalidArgument, "Please supply valid id")
		}
		packageFeatureModel.Pb.Id = in.GetId()
	}

	err := packageFeatureModel.Get(ctx, u.Db)
	if err != nil {
		return &output, err
	}

	return &packageFeatureModel.Pb, nil
}

// List PackageFeature
func (u *PackageFeature) List(in *users.Empty, stream users.PackageFeatureService_ListServer) error {
	ctx := stream.Context()
	rows, err := u.Db.QueryContext(ctx, `SELECT id, name from package_features`)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		err := contextError(ctx)
		if err != nil {
			return err
		}

		var pbPackageFeature users.PackageOfFeature
		var name string
		err = rows.Scan(&pbPackageFeature.Id, &name)
		if err != nil {
			return err
		}

		if value, ok := users.EnumPackageOfFeature_value[name]; ok {
			pbPackageFeature.Name = users.EnumPackageOfFeature(value)
		}

		err = stream.Send(&users.ListPackageFeatureResponse{PackageOfFeature: &pbPackageFeature})
		if err != nil {
			return status.Errorf(codes.Unknown, "cannot send stream response: %v", err)
		}
	}
	return nil
}
