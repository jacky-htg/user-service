package model

import (
	"context"
	"database/sql"

	"github.com/jacky-htg/erp-proto/go/pb/users"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Feature struct
type Feature struct {
	Pb        users.Feature
	PackageID string
}

// FeaturePackage struct
type FeaturePackage struct {
	Pb users.PackageOfFeature
}

// GetByName func
func (u *FeaturePackage) GetByName(ctx context.Context, db *sql.DB) error {
	query := `SELECT id, name FROM package_features WHERE name = $1`
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare statement: %v", err)
	}
	defer stmt.Close()
	var name string
	err = stmt.QueryRowContext(ctx, u.Pb.GetName().String()).Scan(&u.Pb.Id, &name)

	if err == sql.ErrNoRows {
		return status.Errorf(codes.NotFound, "Query Raw: %v", err)
	}

	if err != nil {
		return status.Errorf(codes.Internal, "Query Raw: %v", err)
	}

	feature := Feature{PackageID: u.Pb.GetId()}
	u.Pb.Features, err = feature.GetByPackage(ctx, db)
	if err != nil {
		return err
	}

	return nil
}

// Get func
func (u *FeaturePackage) Get(ctx context.Context, db *sql.DB) error {
	query := `SELECT id, name FROM package_features WHERE id = $1`
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare statement get package feature: %v", err)
	}
	defer stmt.Close()
	var name string
	err = stmt.QueryRowContext(ctx, u.Pb.GetId()).Scan(&u.Pb.Id, &name)

	if err == sql.ErrNoRows {
		return status.Errorf(codes.NotFound, "Query Raw get package feature: %v", err)
	}

	if err != nil {
		return status.Errorf(codes.Internal, "Query Raw get package feature: %v", err)
	}

	if value, ok := users.EnumPackageOfFeature_value[name]; ok {
		u.Pb.Name = users.EnumPackageOfFeature(value)
	}

	if u.Pb.Name == users.EnumPackageOfFeature(0) {
		feature := Feature{}
		u.Pb.Features, err = feature.GetAll(ctx, db)
		if err != nil {
			return err
		}
	} else {
		feature := Feature{PackageID: u.Pb.GetId()}
		u.Pb.Features, err = feature.GetByPackage(ctx, db)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetByPackage feature
func (u *Feature) GetByPackage(ctx context.Context, db *sql.DB) ([]*users.Feature, error) {
	var list []*users.Feature
	query := `
		SELECT features.id, features.name 
		FROM features
		JOIN features_package_features ON features.id = features_package_features.feature_id
		WHERE features_package_features.package_feature_id = $1
	`
	rows, err := db.QueryContext(ctx, query, u.PackageID)
	if err != nil {
		return list, status.Error(codes.Internal, err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var pbFeature users.Feature
		err = rows.Scan(&pbFeature.Id, &pbFeature.Name)
		if err != nil {
			return list, status.Error(codes.Internal, err.Error())
		}

		list = append(list, &pbFeature)
	}

	return list, nil
}

// Get Feature
func (u *Feature) Get(ctx context.Context, db *sql.DB) error {
	query := `SELECT id, name FROM features WHERE id = $1`
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare statement: %v", err)
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, u.Pb.GetName()).Scan(&u.Pb.Id, &u.Pb.Name)

	if err == sql.ErrNoRows {
		return status.Errorf(codes.NotFound, "Query Raw: %v", err)
	}

	if err != nil {
		return status.Errorf(codes.Internal, "Query Raw: %v", err)
	}

	return nil
}

// GetAll Features
func (u *Feature) GetAll(ctx context.Context, db *sql.DB) ([]*users.Feature, error) {
	var list []*users.Feature
	rows, err := db.QueryContext(ctx, `SELECT id, name from features`)
	if err != nil {
		return list, status.Error(codes.Internal, err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var pbFeature users.Feature
		err = rows.Scan(&pbFeature.Id, &pbFeature.Name)
		if err != nil {
			return list, err
		}

		list = append(list, &pbFeature)

	}

	return list, nil
}
