package model

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jacky-htg/erp-pkg/app"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AccessGroup struct
type AccessGroup struct {
	AccessID string
	GroupID  string
}

// Get Access
func (u *AccessGroup) Get(ctx context.Context, db *sql.DB) error {
	var isExist bool
	err := db.QueryRowContext(ctx, `SELECT true FROM access_groups WHERE group_id = $1 AND access_id = $2`, u.GroupID, u.AccessID).Scan(&isExist)

	if err == sql.ErrNoRows {
		return status.Errorf(codes.NotFound, "Query Raw: %v", err)
	}

	if err != nil {
		return status.Errorf(codes.Internal, "Query Raw: %v", err)
	}

	return nil
}

// Grant Access
func (u *AccessGroup) Grant(ctx context.Context, db *sql.DB) error {
	now := time.Now().UTC()
	query := `
		INSERT INTO access_groups (id, group_id, access_id, created_at, created_by, updated_at, updated_by) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare grant access: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		uuid.New().String(),
		u.GroupID,
		u.AccessID,
		now,
		ctx.Value(app.Ctx("userID")).(string),
		now,
		ctx.Value(app.Ctx("userID")).(string),
	)
	if err != nil {
		return status.Errorf(codes.Internal, "Exec grant access: %v", err)
	}

	return nil
}

// Revoke Access
func (u *AccessGroup) Revoke(ctx context.Context, db *sql.DB) error {
	stmt, err := db.PrepareContext(ctx, `DELETE FROM access_groups WHERE group_id = $1 AND access_id = $2`)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare revoke access: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, u.GroupID, u.AccessID)
	if err != nil {
		return status.Errorf(codes.Internal, "Exec rebvoke access: %v", err)
	}

	return nil
}
