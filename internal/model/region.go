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
		return status.Errorf(codes.Internal, "Prepare statement: %v", err)
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, u.Pb.GetId()).Scan(&u.Pb.Id, &u.Pb.CompanyId, &u.Pb.Name, &u.Pb.Code)

	if err == sql.ErrNoRows {
		return status.Errorf(codes.NotFound, "Query Raw: %v", err)
	}

	if err != nil {
		return status.Errorf(codes.Internal, "Query Raw: %v", err)
	}

	return nil
}

// GetByCode Region
func (u *Region) GetByCode(ctx context.Context, db *sql.DB) error {
	query := `SELECT id, company_id, name, code FROM regions WHERE code = $1`
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare statement: %v", err)
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, u.Pb.GetCode()).Scan(&u.Pb.Id, &u.Pb.CompanyId, &u.Pb.Name, &u.Pb.Code)

	if err == sql.ErrNoRows {
		return status.Errorf(codes.NotFound, "Query Raw: %v", err)
	}

	if err != nil {
		return status.Errorf(codes.Internal, "Query Raw: %v", err)
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

func (u *Region) regionBranches(ctx context.Context, tx *sql.Tx, branches []*users.Branch) error {
	for _, branch := range branches {
		query := `INSERT INTO branches_regions (id, region_id, branch_id, created_by, updated_by)
		VALUES ($1, $2, $3, $4, $4)`

		stmt, err := tx.PrepareContext(ctx, query)
		if err != nil {
			return status.Errorf(codes.Internal, "Prepare insert branches_regions: %v", err)
		}
		defer stmt.Close()

		_, err = stmt.ExecContext(ctx,
			uuid.New().String(),
			u.Pb.GetId(),
			branch.GetId(),
			ctx.Value(app.Ctx("userID")).(string),
		)
		if err != nil {
			return status.Errorf(codes.Internal, "exec insert branches_regions: %v", err)
		}
	}

	return nil
}
