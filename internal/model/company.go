package model

import (
	"context"
	"database/sql"
	"time"
	"user-service/internal/pkg/app"
	"user-service/pb/users"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Company struct
type Company struct {
	Pb             users.Company
	UpdateFeatures bool
}

// CompanyRegister struct
type CompanyRegister struct {
	Pb       users.CompanyRegistration
	Password string
}

// Registration Company
func (u *CompanyRegister) Registration(ctx context.Context, db *sql.DB, tx *sql.Tx) error {
	u.Pb.GetCompany().Id = uuid.New().String()
	u.Pb.GetUser().Id = uuid.New().String()
	groupID := uuid.New().String()

	// Create Company
	{
		query := `
		INSERT INTO companies (id, name, code, address, city, province, phone, pic, pic_phone, package_of_feature_id, updated_by, npwp, logo)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		`

		stmt, err := tx.PrepareContext(ctx, query)
		if err != nil {
			return err
		}
		defer stmt.Close()

		packageFeature := FeaturePackage{}
		packageFeature.Pb.Name = u.Pb.GetCompany().GetPackageOfFeature()
		err = packageFeature.GetByName(ctx, db)
		if err != nil {
			return err
		}

		_, err = stmt.ExecContext(ctx,
			u.Pb.GetCompany().GetId(),
			u.Pb.GetCompany().GetName(),
			u.Pb.GetCompany().GetCode(),
			u.Pb.GetCompany().GetAddress(),
			u.Pb.GetCompany().GetCity(),
			u.Pb.GetCompany().GetProvince(),
			u.Pb.GetCompany().GetPhone(),
			u.Pb.GetCompany().GetPic(),
			u.Pb.GetCompany().GetPicPhone(),
			packageFeature.Pb.GetId(),
			u.Pb.GetUser().GetId(),
			u.Pb.GetCompany().GetNpwp(),
			u.Pb.GetCompany().GetLogo(),
		)
		if err != nil {
			return err
		}
	}

	// Create Group
	{
		query := `
			INSERT INTO groups (id, company_id, is_mutable, name, created_by, updated_by)
			VALUES ($1, $2, false, 'Super Admin', $3, $3)
		`

		stmt, err := tx.PrepareContext(ctx, query)
		if err != nil {
			return status.Errorf(codes.Internal, "prepare insert group: %v", err)
		}
		defer stmt.Close()

		_, err = stmt.ExecContext(ctx, groupID, u.Pb.GetCompany().GetId(), u.Pb.GetUser().GetId())
		if err != nil {
			return status.Errorf(codes.Internal, "exec insert group: %v", err)
		}
	}

	// Create User
	{
		query := `
			INSERT INTO users (id, company_id, group_id, username, name, email, password)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`
		pass, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return status.Errorf(codes.Internal, "hash password: %v", err)
		}

		stmt, err := tx.PrepareContext(ctx, query)
		if err != nil {
			return status.Errorf(codes.Internal, "Prepare insert user: %v", err)
		}
		defer stmt.Close()

		_, err = stmt.ExecContext(ctx,
			u.Pb.GetUser().GetId(),
			u.Pb.GetCompany().GetId(),
			groupID,
			u.Pb.GetUser().GetUsername(),
			u.Pb.GetUser().GetName(),
			u.Pb.GetUser().GetEmail(),
			string(pass),
		)
		if err != nil {
			return status.Errorf(codes.Internal, "Exec insert: %v", err)
		}
	}

	// grant access
	{
		var accessModel Access
		err := accessModel.GetRoot(ctx, tx, false)
		if err != nil {
			return err
		}

		query := `
			INSERT INTO access_groups (id, group_id, access_id, created_by, updated_by)
			VALUES ($1, $2, $3, $4, $4)
		`

		stmt, err := tx.PrepareContext(ctx, query)
		if err != nil {
			return err
		}
		defer stmt.Close()

		_, err = stmt.ExecContext(ctx, uuid.New().String(), groupID, accessModel.Pb.GetId(), u.Pb.GetUser().GetId())
		if err != nil {
			return err
		}
	}

	if len(u.Pb.GetCompany().GetFeatures()) > 0 {
		modelCompany := Company{}
		err := modelCompany.featureSetting(ctx, tx, u.Pb.GetCompany().GetFeatures())
		if err != nil {
			return err
		}
	}

	return nil
}

// GetByCode company
func (u *Company) GetByCode(ctx context.Context, db *sql.DB) error {
	query := `SELECT id, name, code, address, city, province, npwp, phone, pic, pic_phone, logo, package_of_feature_id 
	FROM companies WHERE code = $1`
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare statement: %v", err)
	}
	defer stmt.Close()

	var npwp, logo sql.NullString
	var enumPackage string
	err = stmt.QueryRowContext(ctx, u.Pb.GetCode()).Scan(
		&u.Pb.Id, &u.Pb.Name, &u.Pb.Code,
		&u.Pb.Address, &u.Pb.City, &u.Pb.Province,
		&npwp, &u.Pb.Phone, &u.Pb.Pic, &u.Pb.PicPhone, &logo, &enumPackage)

	if err == sql.ErrNoRows {
		return status.Errorf(codes.NotFound, "Query Raw: %v", err)
	}

	if err != nil {
		return status.Errorf(codes.Internal, "Query Raw: %v", err)
	}

	u.Pb.Npwp = npwp.String
	u.Pb.Logo = logo.String
	if value, ok := users.EnumPackageOfFeature_value[enumPackage]; ok {
		u.Pb.PackageOfFeature = users.EnumPackageOfFeature(value)
	}

	return nil
}

// Get company
func (u *Company) Get(ctx context.Context, db *sql.DB) error {
	query := `SELECT id, name, code, address, city, province, npwp, phone, pic, pic_phone, logo, package_of_feature_id 
	FROM companies WHERE id = $1`
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare statement: %v", err)
	}
	defer stmt.Close()

	var npwp, logo sql.NullString
	var enumPackage string
	err = stmt.QueryRowContext(ctx, u.Pb.GetId()).Scan(
		&u.Pb.Id, &u.Pb.Name, &u.Pb.Code,
		&u.Pb.Address, &u.Pb.City, &u.Pb.Province,
		&npwp, &u.Pb.Phone, &u.Pb.Pic, &u.Pb.PicPhone, &logo, &enumPackage)

	if err == sql.ErrNoRows {
		return status.Errorf(codes.NotFound, "Query Raw: %v", err)
	}

	if err != nil {
		return status.Errorf(codes.Internal, "Query Raw: %v", err)
	}

	u.Pb.Npwp = npwp.String
	u.Pb.Logo = logo.String
	if value, ok := users.EnumPackageOfFeature_value[enumPackage]; ok {
		u.Pb.PackageOfFeature = users.EnumPackageOfFeature(value)
	}

	return nil
}

// Update Company
func (u *Company) Update(ctx context.Context, db *sql.DB, tx *sql.Tx) error {
	query := `
		UPDATE companies SET 
		name = $1,
		address = $2, 
		city = $3, 
		province = $4, 
		npwp = $5, 
		phone = $6, 
		pic = $7, 
		pic_phone = $8, 
		logo = $9, 
		package_of_feature_id = $10,
		updated_by = $11,
		updated_at = $12
		WHERE id = $13
	`
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare update: %v", err)
	}
	defer stmt.Close()

	packageFeature := FeaturePackage{}
	packageFeature.Pb.Name = u.Pb.GetPackageOfFeature()
	err = packageFeature.GetByName(ctx, db)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx,
		u.Pb.GetName(),
		u.Pb.GetAddress(),
		u.Pb.GetCity(),
		u.Pb.GetProvince(),
		u.Pb.GetNpwp(),
		u.Pb.GetPhone(),
		u.Pb.GetPic(),
		u.Pb.GetPicPhone(),
		u.Pb.GetLogo(),
		packageFeature.Pb.GetId(),
		ctx.Value(app.Ctx("userID")).(string),
		time.Now().UTC(),
		u.Pb.GetId(),
	)
	if err != nil {
		return status.Errorf(codes.Internal, "Exec update: %v", err)
	}

	if u.UpdateFeatures && len(u.Pb.GetFeatures()) > 0 {
		err = u.featureSetting(ctx, tx, u.Pb.GetFeatures())
		if err != nil {
			return err
		}
	}

	return nil
}

func (u *Company) featureSetting(ctx context.Context, tx *sql.Tx, features []*users.Feature) error {
	for _, feature := range features {
		query := `INSERT INTO companies_features (id, company_id, feature_id, created_by, updated_by)
		VALUES ($1, $2, $3, $4, $4)`

		stmt, err := tx.PrepareContext(ctx, query)
		if err != nil {
			return status.Errorf(codes.Internal, "Prepare insert companies_features: %v", err)
		}
		defer stmt.Close()

		_, err = stmt.ExecContext(ctx,
			uuid.New().String(),
			u.Pb.GetId(),
			feature.GetId(),
			ctx.Value(app.Ctx("userID")).(string),
		)
		if err != nil {
			return status.Errorf(codes.Internal, "exec insert companies_features: %v", err)
		}
	}

	return nil
}
