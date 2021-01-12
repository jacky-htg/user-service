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

// Employee struct
type Employee struct {
	Pb users.Employee
}

// GetByCode func
func (u *Employee) GetByCode(ctx context.Context, db *sql.DB) error {
	query := `
	SELECT employees.id, employees.name, employees.code, employees.address, 
		employees.city, employees.province, employees.jabatan,
		users.id, users.company_id, users.region_id, users.branch_id, users.name, users.email 
	FROM employees 
	JOIN users ON employees.user_id = users.id
	WHERE employees.code = $1`
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare statement get employee by code: %v", err)
	}
	defer stmt.Close()

	var pbUser users.User
	err = stmt.QueryRowContext(ctx, u.Pb.GetCode()).Scan(
		&u.Pb.Id, &u.Pb.Name, &u.Pb.Code, &u.Pb.Address, &u.Pb.City, &u.Pb.Province, &u.Pb.Jabatan,
		&pbUser.Id, &pbUser.CompanyId, &pbUser.RegionId, &pbUser.BranchId, &pbUser.Name, &pbUser.Email,
	)

	if err == sql.ErrNoRows {
		return status.Errorf(codes.NotFound, "Query Raw get employee by code: %v", err)
	}

	if err != nil {
		return status.Errorf(codes.Internal, "Query Raw get employee by code: %v", err)
	}

	u.Pb.User = &pbUser

	return nil
}

// Create Employee
func (u *Employee) Create(ctx context.Context, db *sql.DB) error {
	u.Pb.Id = uuid.New().String()
	now := time.Now().UTC()
	u.Pb.CreatedBy = ctx.Value(app.Ctx("userID")).(string)
	u.Pb.UpdatedBy = ctx.Value(app.Ctx("userID")).(string)

	query := `
		INSERT INTO employees (id, user_id, name, code, address, city, province, jabatan, created_at, created_by, updated_at, updated_by) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare insert employee: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		u.Pb.GetId(),
		u.Pb.GetUser().GetId(),
		u.Pb.GetName(),
		u.Pb.GetCode(),
		u.Pb.GetAddress(),
		u.Pb.GetCity(),
		u.Pb.GetProvince(),
		u.Pb.GetJabatan(),
		now,
		u.Pb.GetCreatedBy(),
		now,
		u.Pb.GetUpdatedBy(),
	)
	if err != nil {
		return status.Errorf(codes.Internal, "Exec insert employee: %v", err)
	}

	u.Pb.CreatedAt, err = ptypes.TimestampProto(now)
	if err != nil {
		return status.Errorf(codes.Internal, "convert created by: %v", err)
	}

	u.Pb.UpdatedAt = u.Pb.CreatedAt

	return nil
}
