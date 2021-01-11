package model

import (
	"context"
	"database/sql"
	"time"
	"user-service/internal/pkg/app"
	users "user-service/pb"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Branch model
type Branch struct {
	Pb           users.Branch
	UpdateRegion bool
}

// Get func
func (u *Branch) Get(ctx context.Context, db *sql.DB) error {
	query := `SELECT branches.id, branches.company_id, regions.id, branches.name, branches.code, branches.address, 
	branches.city, branches.province, branches.npwp, branches.phone, branches.pic, branches.pic_phone 
	FROM branches 
	JOIN branches_regions ON branches.id = branches_regions.branch_id
	JOIN regions ON branches_regions.region_id = regions.id
	WHERE branches.id = $1`
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare statement Get branch: %v", err)
	}
	defer stmt.Close()

	var npwp sql.NullString
	err = stmt.QueryRowContext(ctx, u.Pb.GetId()).Scan(
		&u.Pb.Id, &u.Pb.CompanyId, &u.Pb.RegionId, &u.Pb.Name, &u.Pb.Code, &u.Pb.Address, &u.Pb.City, &u.Pb.Province,
		&npwp, &u.Pb.Phone, &u.Pb.Pic, &u.Pb.PicPhone,
	)

	if err == sql.ErrNoRows {
		return status.Errorf(codes.NotFound, "Query Raw get branch: %v", err)
	}

	if err != nil {
		return status.Errorf(codes.Internal, "Query Raw get branch: %v", err)
	}

	u.Pb.Npwp = npwp.String

	return nil
}

// GetByCode func
func (u *Branch) GetByCode(ctx context.Context, db *sql.DB) error {
	query := `SELECT branches.id, branches.company_id, regions.id, branches.name, branches.code, branches.address, 
	branches.city, branches.province, branches.npwp, branches.phone, branches.pic, branches.pic_phone 
	FROM branches 
	JOIN branches_regions ON branches.id = branches_regions.branch_id
	JOIN regions ON branches_regions.region_id = regions.id
	WHERE branches.code = $1`
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare statement get branch by code: %v", err)
	}
	defer stmt.Close()

	var npwp sql.NullString
	err = stmt.QueryRowContext(ctx, u.Pb.GetCode()).Scan(
		&u.Pb.Id, &u.Pb.CompanyId, &u.Pb.RegionId, &u.Pb.Name, &u.Pb.Code, &u.Pb.Address, &u.Pb.City, &u.Pb.Province,
		&npwp, &u.Pb.Phone, &u.Pb.Pic, &u.Pb.PicPhone,
	)

	if err == sql.ErrNoRows {
		return status.Errorf(codes.NotFound, "Query Raw get branch by code: %v", err)
	}

	if err != nil {
		return status.Errorf(codes.Internal, "Query Raw get branch by code: %v", err)
	}

	u.Pb.Npwp = npwp.String

	return nil
}

// Create Branch
func (u *Branch) Create(ctx context.Context, db *sql.DB, tx *sql.Tx) error {
	u.Pb.Id = uuid.New().String()
	now := time.Now().UTC()
	u.Pb.CreatedBy = ctx.Value(app.Ctx("userID")).(string)
	u.Pb.UpdatedBy = ctx.Value(app.Ctx("userID")).(string)

	query := `
		INSERT INTO branches (id, company_id, name, code, address, city, province, npwp, phone, pic, pic_phone, created_at, created_by, updated_at, updated_by) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	`
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare insert region: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		u.Pb.GetId(),
		u.Pb.GetCompanyId(),
		u.Pb.GetName(),
		u.Pb.GetCode(),
		u.Pb.GetAddress(),
		u.Pb.GetCity(),
		u.Pb.GetProvince(),
		u.Pb.GetNpwp(),
		u.Pb.GetPhone(),
		u.Pb.GetPic(),
		u.Pb.GetPicPhone(),
		now,
		u.Pb.GetCreatedBy(),
		now,
		u.Pb.GetUpdatedBy(),
	)
	if err != nil {
		return status.Errorf(codes.Internal, "Exec insert branch: %v", err)
	}

	branchesRegion := BranchesRegion{RegionID: u.Pb.GetRegionId(), BranchID: u.Pb.GetId()}
	err = branchesRegion.Create(ctx, tx)
	if err != nil {
		return err
	}

	return nil
}

// Update Branch
func (u *Branch) Update(ctx context.Context, db *sql.DB, tx *sql.Tx) error {
	query := `
		UPDATE branches SET 
		name = $1,
		address = $2, 
		city = $3, 
		province = $4, 
		npwp = $5, 
		phone = $6, 
		pic = $7, 
		pic_phone = $8,
		updated_by = $9,
		updated_at = $10
		WHERE id = $11
	`
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare update branch: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		u.Pb.GetName(),
		u.Pb.GetAddress(),
		u.Pb.GetCity(),
		u.Pb.GetProvince(),
		u.Pb.GetNpwp(),
		u.Pb.GetPhone(),
		u.Pb.GetPic(),
		u.Pb.GetPicPhone(),
		ctx.Value(app.Ctx("userID")).(string),
		time.Now().UTC(),
		u.Pb.GetId(),
	)
	if err != nil {
		return status.Errorf(codes.Internal, "Exec update branch: %v", err)
	}

	if u.UpdateRegion && len(u.Pb.GetRegionId()) > 0 {
		// delete current branchesRegion
		{
			branchesRegion := BranchesRegion{BranchID: u.Pb.GetId()}
			err = branchesRegion.DeleteAll(ctx, tx)
			if err != nil {
				return err
			}
		}

		// create new branchesRegion
		{
			branchesRegion := BranchesRegion{RegionID: u.Pb.GetRegionId(), BranchID: u.Pb.GetId()}
			err = branchesRegion.Create(ctx, tx)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
