package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
	"user-service/internal/pkg/app"
	"user-service/pb/users"

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
	WHERE branches.company_id = $1 AND branches.code = $2`
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare statement get branch by code: %v", err)
	}
	defer stmt.Close()

	var npwp sql.NullString
	err = stmt.QueryRowContext(ctx, ctx.Value(app.Ctx("companyID")).(string), u.Pb.GetCode()).Scan(
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
		return status.Errorf(codes.Internal, "Prepare insert branch: %v", err)
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

	u.Pb.CreatedAt = now.String()
	u.Pb.UpdatedAt = u.Pb.CreatedAt

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

// Delete branch
func (u *Branch) Delete(ctx context.Context, db *sql.DB) error {
	stmt, err := db.PrepareContext(ctx, `DELETE FROM branches WHERE id = $1`)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare delete branch: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, u.Pb.GetId())
	if err != nil {
		return status.Errorf(codes.Internal, "Exec delete branch: %v", err)
	}

	return nil
}

// ListQuery builder
func (u *Branch) ListQuery(ctx context.Context, db *sql.DB, in *users.ListBranchRequest, branchID string) (string, []interface{}, *users.BranchPaginationResponse, error) {
	var paginationResponse users.BranchPaginationResponse
	query := `SELECT branches.id, branches.company_id, regions.id, regions.name region_name, branches.name, branches.code, branches.address, 
	branches.city, branches.province, branches.npwp, branches.phone, branches.pic, branches.pic_phone 
	FROM branches 
	JOIN branches_regions ON branches.id = branches_regions.branch_id
	JOIN regions ON branches_regions.region_id = regions.id`
	where := []string{"branches.company_id = $1"}
	paramQueries := []interface{}{ctx.Value(app.Ctx("companyID")).(string)}

	if len(branchID) > 0 {
		paramQueries = append(paramQueries, branchID)
		where = append(where, fmt.Sprintf(`branches.id = $%d`, len(paramQueries)))
	}

	if len(in.GetRegionId()) > 0 {
		paramQueries = append(paramQueries, in.GetRegionId())
		where = append(where, fmt.Sprintf(`regions.id = $%d`, len(paramQueries)))
	}

	if len(in.GetPagination().GetSearch()) > 0 {
		paramQueries = append(paramQueries, in.GetPagination().GetSearch()+"%")
		where = append(where, fmt.Sprintf(`(branches.name ILIKE $%d OR branches.code ILIKE $%d 
			OR branches.address ILIKE $%d OR branches.city ILIKE $%d OR branches.province ILIKE $%d 
			OR branches.npwp ILIKE $%d OR branches.phone ILIKE $%d OR branches.pic ILIKE $%d
			OR branches.pic_phone ILIKE $%d OR regions.name ILIKE $%d)`,
			len(paramQueries), len(paramQueries), len(paramQueries), len(paramQueries), len(paramQueries),
			len(paramQueries), len(paramQueries), len(paramQueries), len(paramQueries), len(paramQueries)))
	}

	{
		qCount := `SELECT COUNT(*) FROM branches 
		JOIN branches_regions ON branches.id = branches_regions.branch_id
		JOIN regions ON branches_regions.region_id = regions.id`
		if len(where) > 0 {
			qCount += " WHERE " + strings.Join(where, " AND ")
		}
		var count int
		err := db.QueryRowContext(ctx, qCount, paramQueries...).Scan(&count)
		if err != nil && err != sql.ErrNoRows {
			return query, paramQueries, &paginationResponse, status.Error(codes.Internal, err.Error())
		}

		paginationResponse.Count = uint32(count)
	}

	if len(where) > 0 {
		query += ` WHERE ` + strings.Join(where, " AND ")
	}

	if len(in.GetPagination().GetOrderBy()) == 0 || !(in.GetPagination().GetOrderBy() == "branches.name" ||
		in.GetPagination().GetOrderBy() == "branches.code") {
		if in.GetPagination() == nil {
			in.Pagination = &users.Pagination{OrderBy: "branches.created_at"}
		} else {
			in.GetPagination().OrderBy = "branches.created_at"
		}
	}

	query += ` ORDER BY ` + in.GetPagination().GetOrderBy() + ` ` + in.GetPagination().GetSort().String()

	if in.GetPagination().GetLimit() > 0 {
		query += fmt.Sprintf(` LIMIT $%d OFFSET $%d`, (len(paramQueries) + 1), (len(paramQueries) + 2))
		paramQueries = append(paramQueries, in.GetPagination().GetLimit(), in.GetPagination().GetOffset())
	}

	return query, paramQueries, &paginationResponse, nil
}
