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

// Region model
type Region struct {
	Pb             users.Region
	UpdateBranches bool
}

// Get func
func (u *Region) Get(ctx context.Context, db *sql.DB) error {
	query := `SELECT id, company_id, name, code FROM regions WHERE id = $1`
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare statement get region: %v", err)
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, u.Pb.GetId()).Scan(&u.Pb.Id, &u.Pb.CompanyId, &u.Pb.Name, &u.Pb.Code)

	if err == sql.ErrNoRows {
		return status.Errorf(codes.NotFound, "Query Raw get region: %v", err)
	}

	if err != nil {
		return status.Errorf(codes.Internal, "Query Raw get region: %v", err)
	}

	return nil
}

// GetByCode Region
func (u *Region) GetByCode(ctx context.Context, db *sql.DB) error {
	query := `SELECT id, company_id, name, code FROM regions WHERE company_id = $1 AND code = $2`
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare statement get region by code: %v", err)
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, ctx.Value(app.Ctx("companyID")).(string), u.Pb.GetCode()).Scan(&u.Pb.Id, &u.Pb.CompanyId, &u.Pb.Name, &u.Pb.Code)

	if err == sql.ErrNoRows {
		return status.Errorf(codes.NotFound, "Query Raw get region by code: %v", err)
	}

	if err != nil {
		return status.Errorf(codes.Internal, "Query Raw get region by code: %v", err)
	}

	return nil
}

// Create Region
func (u *Region) Create(ctx context.Context, db *sql.DB, tx *sql.Tx) error {
	u.Pb.Id = uuid.New().String()
	query := `
		INSERT INTO regions (id, company_id, name, code, created_by, updated_by) 
		VALUES ($1, $2, $3, $4, $5, $5)
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
		ctx.Value(app.Ctx("userID")).(string),
	)
	if err != nil {
		return status.Errorf(codes.Internal, "Exec insert region: %v", err)
	}

	if len(u.Pb.GetBranches()) > 0 {
		err = u.regionBranches(ctx, tx, u.Pb.GetBranches())
		if err != nil {
			return err
		}
	}

	return nil
}

// Update Region
func (u *Region) Update(ctx context.Context, db *sql.DB, tx *sql.Tx) error {
	query := `
		UPDATE regions SET 
		name = $1,
		updated_by = $2,
		updated_at = $3
		WHERE id = $4
	`
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare update region: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		u.Pb.GetName(),
		ctx.Value(app.Ctx("userID")).(string),
		time.Now().UTC(),
		u.Pb.GetId(),
	)
	if err != nil {
		return status.Errorf(codes.Internal, "Exec update region: %v", err)
	}

	if u.UpdateBranches && len(u.Pb.GetBranches()) > 0 {
		err = u.regionBranches(ctx, tx, u.Pb.GetBranches())
		if err != nil {
			return err
		}
	}

	return nil
}

// Delete region
func (u *Region) Delete(ctx context.Context, db *sql.DB) error {
	stmt, err := db.PrepareContext(ctx, `DELETE FROM regions WHERE id = $1`)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare delete region: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, u.Pb.GetId())
	if err != nil {
		return status.Errorf(codes.Internal, "Exec delete region: %v", err)
	}

	return nil
}

// ListQuery builder
func (u *Region) ListQuery(ctx context.Context, db *sql.DB, in *users.ListRegionRequest, regionID string) (string, []interface{}, *users.RegionPaginationResponse, error) {
	var paginationResponse users.RegionPaginationResponse
	query := `SELECT id, company_id, name, code FROM regions`
	where := []string{"company_id = $1"}
	paramQueries := []interface{}{ctx.Value(app.Ctx("companyID")).(string)}

	if len(regionID) > 0 {
		paramQueries = append(paramQueries, regionID)
		where = append(where, fmt.Sprintf(`id = $%d`, len(paramQueries)))
	}

	if len(in.GetPagination().GetSearch()) > 0 {
		paramQueries = append(paramQueries, in.GetPagination().GetSearch())
		where = append(where, fmt.Sprintf(`(name ILIKE $%d OR code ILIKE $%d)`,
			len(paramQueries), len(paramQueries)))
	}

	{
		qCount := `SELECT COUNT(*) FROM regions`
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

	if len(in.GetPagination().GetOrderBy()) == 0 || !(in.GetPagination().GetOrderBy() == "name" ||
		in.GetPagination().GetOrderBy() == "code") {
		if in.GetPagination() == nil {
			in.Pagination = &users.Pagination{OrderBy: "created_at"}
		} else {
			in.GetPagination().OrderBy = "created_at"
		}
	}

	query += ` ORDER BY ` + in.GetPagination().GetOrderBy() + ` ` + in.GetPagination().GetSort().String()

	if in.GetPagination().GetLimit() > 0 {
		query += fmt.Sprintf(` LIMIT $%d OFFSET $%d`, (len(paramQueries) + 1), (len(paramQueries) + 2))
		paramQueries = append(paramQueries, in.GetPagination().GetLimit(), in.GetPagination().GetOffset())
	}

	return query, paramQueries, &paginationResponse, nil
}

func (u *Region) regionBranches(ctx context.Context, tx *sql.Tx, branches []*users.Branch) error {
	for _, branch := range branches {
		branchesRegion := BranchesRegion{RegionID: u.Pb.GetId(), BranchID: branch.GetId()}
		err := branchesRegion.Create(ctx, tx)
		if err != nil {
			return err
		}
	}

	return nil
}
