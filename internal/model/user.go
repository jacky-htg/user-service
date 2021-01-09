package model

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"user-service/internal/pkg/app"
	users "user-service/pb"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// User struct
type User struct {
	Pb       users.User
	Password string
}

// GetByUserNamePassword func
func (u *User) GetByUserNamePassword(ctx context.Context, db *sql.DB) error {
	var strPassword, tmpAccess string
	var regionID, branchID sql.NullString
	query := `
		SELECT users.id, users.company_id, users.region_id, users.branch_id, users.name, users.email, users.password,
		groups.id groups_id, groups.name groups_name, 
		json_agg(DISTINCT jsonb_build_object(
			'id', access.id,
			'name', access.name
		)) as access
		FROM users
		JOIN groups ON users.group_id = groups.id
		LEFT JOIN access_groups ON groups.id = access_groups.group_id
		LEFT JOIN access ON access_groups.access_id = access.id
		WHERE users.username = $1
		GROUP BY users.id, users.company_id, users.region_id, users.branch_id, users.name, users.email, users.password, groups_id, groups_name 
	`
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare statement: %v", err)
	}
	defer stmt.Close()

	group := users.Group{}
	err = stmt.QueryRowContext(ctx, u.Pb.GetUsername()).Scan(
		&u.Pb.Id, &u.Pb.CompanyId, &regionID, &branchID, &u.Pb.Name, &u.Pb.Email, &strPassword,
		&group.Id, &group.Name, &tmpAccess)

	if err == sql.ErrNoRows {
		return status.Errorf(codes.NotFound, "Query Raw: %v", err)
	}

	if err != nil {
		return status.Errorf(codes.Internal, "Query Raw: %v", err)
	}

	u.Pb.RegionId = regionID.String
	u.Pb.BranchId = branchID.String

	err = bcrypt.CompareHashAndPassword([]byte(strPassword), []byte(u.Password))
	if err != nil {
		return status.Errorf(codes.NotFound, "Invalid Password: %v", err)
	}

	err = json.Unmarshal([]byte(tmpAccess), &group.Access)
	if err != nil {
		return status.Errorf(codes.Internal, "unmarshal access: %v", err)
	}

	u.Pb.Group = &group

	return nil
}

// Get func
func (u *User) Get(ctx context.Context, db *sql.DB) error {
	var regionID, branchID sql.NullString
	query := `
		SELECT users.id, users.company_id, users.region_id, users.branch_id, users.name, users.email,
		groups.id groups_id, groups.name groups_name
		FROM users
		JOIN groups ON users.group_id = groups.id
		WHERE users.id = $1 
	`
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare statement: %v", err)
	}
	defer stmt.Close()

	group := users.Group{}
	err = stmt.QueryRowContext(ctx, u.Pb.GetId()).Scan(
		&u.Pb.Id, &u.Pb.CompanyId, &regionID, &branchID, &u.Pb.Name, &u.Pb.Email, &group.Id, &group.Name)

	if err == sql.ErrNoRows {
		return status.Errorf(codes.NotFound, "Query Raw: %v", err)
	}

	if err != nil {
		return status.Errorf(codes.Internal, "Query Raw: %v", err)
	}

	u.Pb.RegionId = regionID.String
	u.Pb.BranchId = branchID.String
	u.Pb.Group = &group

	return nil
}

// GetByPassword func
func (u *User) GetByPassword(ctx context.Context, db *sql.DB) error {
	var regionID, branchID sql.NullString
	var strPassword string
	query := `
		SELECT users.id, users.company_id, users.region_id, users.branch_id, users.name, users.email, users.password,
		groups.id groups_id, groups.name groups_name
		FROM users
		JOIN groups ON users.group_id = groups.id
		WHERE users.id = $1 
	`
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare statement: %v", err)
	}
	defer stmt.Close()

	group := users.Group{}
	err = stmt.QueryRowContext(ctx, u.Pb.GetId()).Scan(
		&u.Pb.Id, &u.Pb.CompanyId, &regionID, &branchID, &u.Pb.Name, &u.Pb.Email, &strPassword, &group.Id, &group.Name)

	if err == sql.ErrNoRows {
		return status.Errorf(codes.NotFound, "Query Raw: %v", err)
	}

	if err != nil {
		return status.Errorf(codes.Internal, "Query Raw: %v", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(strPassword), []byte(u.Password))
	if err != nil {
		return status.Errorf(codes.NotFound, "Invalid Password: %v", err)
	}

	u.Pb.RegionId = regionID.String
	u.Pb.BranchId = branchID.String
	u.Pb.Group = &group

	return nil
}

// GetByEmail func
func (u *User) GetByEmail(ctx context.Context, db *sql.DB) error {
	var regionID, branchID sql.NullString
	query := `
		SELECT users.id, users.company_id, users.region_id, users.branch_id, users.name, users.email,
		groups.id groups_id, groups.name groups_name
		FROM users
		JOIN groups ON users.group_id = groups.id
		WHERE users.email = $1 
	`
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare statement: %v", err)
	}
	defer stmt.Close()

	group := users.Group{}
	err = stmt.QueryRowContext(ctx, u.Pb.GetEmail()).Scan(
		&u.Pb.Id, &u.Pb.CompanyId, &regionID, &branchID, &u.Pb.Name, &u.Pb.Email, &group.Id, &group.Name)

	if err == sql.ErrNoRows {
		return status.Errorf(codes.NotFound, "Query Raw: %v", err)
	}

	if err != nil {
		return status.Errorf(codes.Internal, "Query Raw: %v", err)
	}

	u.Pb.RegionId = regionID.String
	u.Pb.BranchId = branchID.String
	u.Pb.Group = &group

	return nil
}

// GetByUsername func
func (u *User) GetByUsername(ctx context.Context, db *sql.DB) error {
	var regionID, branchID sql.NullString
	query := `
		SELECT users.id, users.company_id, users.region_id, users.branch_id, users.name, users.email,
		groups.id groups_id, groups.name groups_name
		FROM users
		JOIN groups ON users.group_id = groups.id
		WHERE users.username = $1 
	`
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare statement: %v", err)
	}
	defer stmt.Close()

	group := users.Group{}
	err = stmt.QueryRowContext(ctx, u.Pb.GetUsername()).Scan(
		&u.Pb.Id, &u.Pb.CompanyId, &regionID, &branchID, &u.Pb.Name, &u.Pb.Email, &group.Id, &group.Name)

	if err == sql.ErrNoRows {
		return status.Errorf(codes.NotFound, "Query Raw: %v", err)
	}

	if err != nil {
		return status.Errorf(codes.Internal, "Query Raw: %v", err)
	}

	u.Pb.RegionId = regionID.String
	u.Pb.BranchId = branchID.String
	u.Pb.Group = &group

	return nil
}

// ChangePassword func
func (u *User) ChangePassword(ctx context.Context, tx *sql.Tx) error {

	stmt, err := tx.PrepareContext(ctx, `UPDATE users SET password = $1 WHERE id = $2`)
	defer stmt.Close()
	if err != nil {
		return status.Errorf(codes.Internal, "prepare update: %v", err)
	}

	pass, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return status.Errorf(codes.Internal, "hash password: %v", err)
	}
	_, err = stmt.ExecContext(ctx, string(pass), u.Pb.GetId())
	if err != nil {
		return status.Errorf(codes.Internal, "exec update: %v", err)
	}

	return nil
}

// IsAuth func
func (u *User) IsAuth(ctx context.Context, db *sql.DB, access string) error {
	var parent *string
	strArr := strings.Split(access, "::")
	if len(strArr) > 1 {
		parent = &strArr[0]
	}
	query := `
		SELECT 1
		FROM users
		JOIN groups ON users.group_id = groups.id
		JOIN access_groups ON groups.id = access_groups.group_id
		JOIN access ON access_groups.access_id = access.id
		WHERE users.id = $1 AND (access.name = 'root' OR access.name = $2 OR access.name = $3)
	`
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare statement: %v", err)
	}
	defer stmt.Close()

	var result bool
	err = stmt.QueryRowContext(ctx, u.Pb.GetId(), parent, access).Scan(&result)

	if err == sql.ErrNoRows {
		return status.Errorf(codes.NotFound, "Query Raw: %v", err)
	}

	if err != nil {
		return status.Errorf(codes.Internal, "Query Raw: %v", err)
	}

	return nil
}

// Create new user
func (u *User) Create(ctx context.Context, db *sql.DB) error {
	query := `
		INSERT INTO users (id, company_id, region_id, branch_id, group_id, username, name, email, password, updated_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	var regionID, branchID *string
	if len(u.Pb.GetRegionId()) > 0 {
		regionID = &u.Pb.RegionId
	}

	if len(u.Pb.GetBranchId()) > 0 {
		branchID = &u.Pb.BranchId
	}

	u.Pb.Id = uuid.New().String()

	pass, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return status.Errorf(codes.Internal, "hash password: %v", err)
	}

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare insert: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		u.Pb.GetId(),
		u.Pb.GetCompanyId(),
		regionID,
		branchID,
		u.Pb.GetGroup().GetId(),
		u.Pb.GetUsername(),
		u.Pb.GetName(),
		u.Pb.GetEmail(),
		string(pass),
		ctx.Value(app.Ctx("userID")).(string),
	)
	if err != nil {
		return status.Errorf(codes.Internal, "Exec insert: %v", err)
	}

	return nil
}

// Update user
func (u *User) Update(ctx context.Context, db *sql.DB) error {
	query := `
		UPDATE users SET 
		region_id = $1, 
		branch_id = $2, 
		group_id = $3, 
		name = $4,
		updated_by = $5,
		updated_at = $6
		WHERE id = $7
	`
	var regionID, branchID *string
	if len(u.Pb.GetRegionId()) > 0 {
		regionID = &u.Pb.RegionId
	}

	if len(u.Pb.GetBranchId()) > 0 {
		branchID = &u.Pb.BranchId
	}

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare update: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		regionID,
		branchID,
		u.Pb.GetGroup().GetId(),
		u.Pb.GetName(),
		ctx.Value(app.Ctx("userID")).(string),
		time.Now().UTC(),
		u.Pb.GetId(),
	)
	if err != nil {
		return status.Errorf(codes.Internal, "Exec update: %v", err)
	}

	return nil
}

// Delete user
func (u *User) Delete(ctx context.Context, db *sql.DB) error {
	stmt, err := db.PrepareContext(ctx, `DELETE FROM users WHERE id = $1`)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare delete: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, u.Pb.GetId())
	if err != nil {
		return status.Errorf(codes.Internal, "Exec delete: %v", err)
	}

	return nil
}

// ListQuery builder
func (u *User) ListQuery(ctx context.Context, db *sql.DB, in *users.ListUserRequest) (string, []interface{}, *users.UserPaginationResponse, error) {
	var paginationResponse users.UserPaginationResponse
	query := `
		SELECT users.id, users.company_id, users.region_id, users.branch_id, users.name, users.email,
			groups.id groups_id, groups.name groups_name
		FROM users
		JOIN groups ON users.group_id = groups.id
	`
	where := []string{}
	paramQueries := []interface{}{}

	if len(in.GetBranchId()) > 0 {
		paramQueries = append(paramQueries, in.GetBranchId)
		where = append(where, fmt.Sprintf(`users.branch_id = $%d`, len(paramQueries)))
	}

	if len(in.GetCompanyId()) > 0 {
		paramQueries = append(paramQueries, in.GetCompanyId)
		where = append(where, fmt.Sprintf(`users.company_id = $%d`, len(paramQueries)))
	}

	if len(in.GetPagination().GetSearch()) > 0 {
		paramQueries = append(paramQueries, in.GetPagination().GetSearch())
		where = append(where, fmt.Sprintf(`(users.name ILIKE $%d OR users.email ILIKE $%d OR groups.name ILIKE $%d)`,
			len(paramQueries), len(paramQueries), len(paramQueries)))
	}

	{
		qCount := `SELECT COUNT(*) FROM users JOIN groups ON users.group_id = groups.id`
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
		query += `WHERE ` + strings.Join(where, " AND ")
	}

	if len(in.GetPagination().GetOrderBy()) == 0 || !(in.GetPagination().GetOrderBy() == "users.name" ||
		in.GetPagination().GetOrderBy() == "users.email" ||
		in.GetPagination().GetOrderBy() == "groups.name") {
		if in.GetPagination() == nil {
			in.Pagination = &users.Pagination{OrderBy: "users.created_at"}
		} else {
			in.GetPagination().OrderBy = "users.created_at"
		}
	}

	query += ` ORDER BY ` + in.GetPagination().GetOrderBy() + ` ` + in.GetPagination().GetSort().String()

	if in.GetPagination().GetLimit() > 0 {
		query += fmt.Sprintf(` LIMIT $%d OFFSET $%d`, (len(paramQueries) + 1), (len(paramQueries) + 2))
		paramQueries = append(paramQueries, in.GetPagination().GetLimit(), in.GetPagination().GetOffset())
	}

	return query, paramQueries, &paginationResponse, nil
}
