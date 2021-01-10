package model

import (
	"context"
	"database/sql"
	users "user-service/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Access struct
type Access struct {
	Pb users.Access
}

// GetRoot Access
func (u *Access) GetRoot(ctx context.Context, tx *sql.Tx) error {
	err := tx.QueryRowContext(ctx, `SELECT id, name FROM access WHERE name = 'root'`).Scan(&u.Pb.Id, &u.Pb.Name)

	if err == sql.ErrNoRows {
		return status.Errorf(codes.NotFound, "Query Raw: %v", err)
	}

	if err != nil {
		return status.Errorf(codes.Internal, "Query Raw: %v", err)
	}

	return nil
}
