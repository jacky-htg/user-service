package model

import (
	"context"
	"database/sql"
	users "user-service/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Region model
type Region struct {
	Pb users.Region
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
