package model

import (
	"context"
	"database/sql"
	"time"
	users "user-service/pb"

	"github.com/golang/protobuf/ptypes"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RequestPassword model
type RequestPassword struct {
	Pb users.RequestPassword
}

// Create func
func (u *RequestPassword) Create(ctx context.Context, db *sql.DB) error {
	u.Pb.Id = uuid.New().String()
	stmt, err := db.PrepareContext(ctx, `INSERT INTO request_passwords (id, user_id, created_at) VALUES ($1, $2, $3)`)
	if err != nil {
		return status.Errorf(codes.Internal, "prepare insert: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, u.Pb.GetId(), u.Pb.GetUserId(), time.Now().UTC())
	if err != nil {
		return status.Errorf(codes.Internal, "exec insert: %v", err)
	}

	return nil
}

// Get func
func (u *RequestPassword) Get(ctx context.Context, db *sql.DB) error {
	query := `SELECT id, user_id, is_used, created_at FROM request_passwords WHERE id = $1`
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare statement: %v", err)
	}
	defer stmt.Close()

	var createdAt time.Time
	err = stmt.QueryRowContext(ctx, u.Pb.GetId()).Scan(&u.Pb.Id, &u.Pb.UserId, &u.Pb.IsUsed, &createdAt)

	if err == sql.ErrNoRows {
		return status.Errorf(codes.NotFound, "Query Raw: %v", err)
	}

	if err != nil {
		return status.Errorf(codes.Internal, "Query Raw: %v", err)
	}

	u.Pb.CreatedAt, err = ptypes.TimestampProto(createdAt)
	if err != nil {
		return status.Errorf(codes.Internal, "convert TimestampProto: %v", err)
	}

	return nil
}

// UpdateIsUsed func
func (u *RequestPassword) UpdateIsUsed(ctx context.Context, tx *sql.Tx) error {

	stmt, err := tx.PrepareContext(ctx, `UPDATE request_passwords SET is_used = $1 WHERE id = $2`)
	defer stmt.Close()
	if err != nil {
		return status.Errorf(codes.Internal, "prepare update: %v", err)
	}

	_, err = stmt.ExecContext(ctx, true, u.Pb.GetId())
	if err != nil {
		return status.Errorf(codes.Internal, "exec update: %v", err)
	}
	return nil
}
