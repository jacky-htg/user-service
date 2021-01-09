package model

import (
	"context"
	"database/sql"
	users "user-service/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Group model
type Group struct {
	Pb users.Group
}

// Get func
func (u *Group) Get(ctx context.Context, db *sql.DB) error {
	query := `SELECT id, company_id, name, is_mutable FROM groups WHERE id = $1`
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare statement: %v", err)
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, u.Pb.GetId()).Scan(&u.Pb.Id, &u.Pb.CompanyId, &u.Pb.Name, &u.Pb.IsMutable)

	if err == sql.ErrNoRows {
		return status.Errorf(codes.NotFound, "Query Raw: %v", err)
	}

	if err != nil {
		return status.Errorf(codes.Internal, "Query Raw: %v", err)
	}

	return nil
}
