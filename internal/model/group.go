package model

import (
	"context"
	"database/sql"
	"time"
	"user-service/internal/pkg/app"
	users "user-service/pb"

	"github.com/golang/protobuf/ptypes"

	"github.com/google/uuid"
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

// Create Group
func (u *Group) Create(ctx context.Context, db *sql.DB) error {
	var err error
	u.Pb.Id = uuid.New().String()
	now := time.Now().UTC()
	u.Pb.CreatedBy = ctx.Value(app.Ctx("userID")).(string)
	u.Pb.UpdatedBy = ctx.Value(app.Ctx("userID")).(string)

	u.Pb.CreatedAt, err = ptypes.TimestampProto(now)
	if err != nil {
		return err
	}

	u.Pb.UpdatedAt = u.Pb.CreatedAt

	query := `
		INSERT INTO groups (id, company_id, name, created_at, created_by, updated_at, updated_by) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare insert group: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		u.Pb.GetId(),
		u.Pb.GetCompanyId(),
		u.Pb.GetName(),
		now,
		u.Pb.GetCreatedBy(),
		now,
		u.Pb.GetUpdatedBy(),
	)
	if err != nil {
		return status.Errorf(codes.Internal, "Exec insert group: %v", err)
	}

	return nil
}

// Update Group
func (u *Group) Update(ctx context.Context, db *sql.DB) error {
	var err error
	now := time.Now().UTC()
	u.Pb.UpdatedBy = ctx.Value(app.Ctx("userID")).(string)

	u.Pb.UpdatedAt, err = ptypes.TimestampProto(now)
	if err != nil {
		return err
	}

	query := `
		UPDATE groups SET
		name = $1, 
		updated_at = $2, 
		updated_by = $3
		WHERE id = $4
	`
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare update group: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		u.Pb.GetName(),
		now,
		u.Pb.GetUpdatedBy(),
		u.Pb.GetId(),
	)
	if err != nil {
		return status.Errorf(codes.Internal, "Exec update group: %v", err)
	}

	return nil
}

// Delete group
func (u *Group) Delete(ctx context.Context, db *sql.DB) error {
	stmt, err := db.PrepareContext(ctx, `DELETE FROM groups WHERE id = $1`)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare delete group: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, u.Pb.GetId())
	if err != nil {
		return status.Errorf(codes.Internal, "Exec delete group: %v", err)
	}

	return nil
}
