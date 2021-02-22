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
	var regionID, branchID sql.NullString
	err = stmt.QueryRowContext(ctx, u.Pb.GetCode()).Scan(
		&u.Pb.Id, &u.Pb.Name, &u.Pb.Code, &u.Pb.Address, &u.Pb.City, &u.Pb.Province, &u.Pb.Jabatan,
		&pbUser.Id, &pbUser.CompanyId, &regionID, &branchID, &pbUser.Name, &pbUser.Email,
	)

	if err == sql.ErrNoRows {
		return status.Errorf(codes.NotFound, "Query Raw get employee by code: %v", err)
	}

	if err != nil {
		return status.Errorf(codes.Internal, "Query Raw get employee by code: %v", err)
	}

	pbUser.RegionId = regionID.String
	pbUser.BranchId = branchID.String
	u.Pb.User = &pbUser

	return nil
}

// Get func
func (u *Employee) Get(ctx context.Context, db *sql.DB) error {
	query := `
	SELECT employees.id, employees.name, employees.code, employees.address, 
		employees.city, employees.province, employees.jabatan,
		users.id, users.company_id, users.region_id, users.branch_id, users.name, users.email 
	FROM employees 
	JOIN users ON employees.user_id = users.id
	WHERE employees.id = $1`
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare statement get employee: %v", err)
	}
	defer stmt.Close()

	var pbUser users.User
	var regionID, branchID sql.NullString
	err = stmt.QueryRowContext(ctx, u.Pb.GetId()).Scan(
		&u.Pb.Id, &u.Pb.Name, &u.Pb.Code, &u.Pb.Address, &u.Pb.City, &u.Pb.Province, &u.Pb.Jabatan,
		&pbUser.Id, &pbUser.CompanyId, &regionID, &branchID, &pbUser.Name, &pbUser.Email,
	)

	if err == sql.ErrNoRows {
		return status.Errorf(codes.NotFound, "Query Raw get employee: %v", err)
	}

	if err != nil {
		return status.Errorf(codes.Internal, "Query Raw get employee: %v", err)
	}

	pbUser.RegionId = regionID.String
	pbUser.BranchId = branchID.String
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

	u.Pb.CreatedAt = now.String()
	u.Pb.UpdatedAt = u.Pb.CreatedAt

	return nil
}

// Update Employee
func (u *Employee) Update(ctx context.Context, db *sql.DB) error {
	now := time.Now().UTC()
	u.Pb.UpdatedBy = ctx.Value(app.Ctx("userID")).(string)
	query := `
		UPDATE employees SET 
		user_id = $1, 
		name = $2, 
		address = $3, 
		city = $4, 
		province = $5, 
		jabatan = $6, 
		updated_at = $7, 
		updated_by = $8
		WHERE id = $9
	`
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare update employee: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		u.Pb.GetUser().GetId(),
		u.Pb.GetName(),
		u.Pb.GetAddress(),
		u.Pb.GetCity(),
		u.Pb.GetProvince(),
		u.Pb.GetJabatan(),
		now,
		u.Pb.GetUpdatedBy(),
		u.Pb.GetId(),
	)
	if err != nil {
		return status.Errorf(codes.Internal, "Exec update employee: %v", err)
	}

	u.Pb.UpdatedAt = now.String()

	return nil
}

// Delete employee
func (u *Employee) Delete(ctx context.Context, db *sql.DB) error {
	stmt, err := db.PrepareContext(ctx, `DELETE FROM employees WHERE id = $1`)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare delete employee: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, u.Pb.GetId())
	if err != nil {
		return status.Errorf(codes.Internal, "Exec delete employee: %v", err)
	}

	return nil
}

// ListQuery builder
func (u *Employee) ListQuery(ctx context.Context, db *sql.DB, in *users.ListEmployeeRequest, userLogin *users.User) (string, []interface{}, *users.EmployeePaginationResponse, error) {
	var paginationResponse users.EmployeePaginationResponse
	query := `
		SELECT employees.id, employees.name, employees.code, employees.address, 
			employees.city, employees.province, employees.jabatan,
			users.id, users.company_id, users.region_id, users.branch_id, users.name, users.email 
		FROM employees 
		JOIN users ON employees.user_id = users.id
	`
	where := []string{"users.company_id = $1"}
	paramQueries := []interface{}{userLogin.GetCompanyId()}

	if len(userLogin.GetRegionId()) > 0 {
		paramQueries = append(paramQueries, userLogin.GetRegionId())
		where = append(where, fmt.Sprintf(`users.region_id = $%d`, len(paramQueries)))
	}

	if len(userLogin.GetBranchId()) > 0 {
		paramQueries = append(paramQueries, userLogin.GetBranchId())
		where = append(where, fmt.Sprintf(`users.branch_id = $%d`, len(paramQueries)))
	}

	if len(in.GetPagination().GetSearch()) > 0 {
		paramQueries = append(paramQueries, in.GetPagination().GetSearch())
		where = append(where, fmt.Sprintf(`(employees.name ILIKE $%d OR employees.code ILIKE $%d 
			OR employees.address ILIKE $%d OR employees.city ILIKE $%d OR employees.province ILIKE $%d 
			OR employees.jabatan ILIKE $%d OR users.name ILIKE $%d OR users.email ILIKE $%d)`,
			len(paramQueries), len(paramQueries), len(paramQueries), len(paramQueries),
			len(paramQueries), len(paramQueries), len(paramQueries), len(paramQueries)))
	}

	{
		qCount := `SELECT COUNT(*) FROM employees JOIN users ON employees.user_id = users.id`
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

	if len(in.GetPagination().GetOrderBy()) == 0 || !(in.GetPagination().GetOrderBy() == "employees.name" ||
		in.GetPagination().GetOrderBy() == "employees.code") {
		if in.GetPagination() == nil {
			in.Pagination = &users.Pagination{OrderBy: "employees.created_at"}
		} else {
			in.GetPagination().OrderBy = "employees.created_at"
		}
	}

	query += ` ORDER BY ` + in.GetPagination().GetOrderBy() + ` ` + in.GetPagination().GetSort().String()

	if in.GetPagination().GetLimit() > 0 {
		query += fmt.Sprintf(` LIMIT $%d OFFSET $%d`, (len(paramQueries) + 1), (len(paramQueries) + 2))
		paramQueries = append(paramQueries, in.GetPagination().GetLimit(), in.GetPagination().GetOffset())
	}

	return query, paramQueries, &paginationResponse, nil
}
