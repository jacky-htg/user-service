package model

import (
	"context"
	"database/sql"
	"user-service/internal/pkg/app"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// BranchesRegion strcut
type BranchesRegion struct {
	ID       string
	BranchID string
	RegionID string
}

// Create branchesRegion
func (u *BranchesRegion) Create(ctx context.Context, tx *sql.Tx) error {
	u.ID = uuid.New().String()

	query := `INSERT INTO branches_regions (id, region_id, branch_id, created_by, updated_by)
		VALUES ($1, $2, $3, $4, $4)`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare insert branches_regions: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		u.ID,
		u.RegionID,
		u.BranchID,
		ctx.Value(app.Ctx("userID")).(string),
	)
	if err != nil {
		return status.Errorf(codes.Internal, "exec insert branches_regions: %v", err)
	}

	return nil
}
