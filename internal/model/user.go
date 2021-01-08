package model

import (
	"context"
	"database/sql"
	"encoding/json"
	"strings"
	users "user-service/pb"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// User struct
type User struct {
	Pb users.User
}

// GetByUserNamePassword func
func (u *User) GetByUserNamePassword(ctx context.Context, db *sql.DB, password string) error {
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

	err = bcrypt.CompareHashAndPassword([]byte(strPassword), []byte(password))
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
func (u *User) GetByPassword(ctx context.Context, db *sql.DB, password string) error {
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

	err = bcrypt.CompareHashAndPassword([]byte(strPassword), []byte(password))
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

// ChangePassword func
func (u *User) ChangePassword(ctx context.Context, tx *sql.Tx, password string) error {

	stmt, err := tx.PrepareContext(ctx, `UPDATE users SET password = $1 WHERE id = $2`)
	defer stmt.Close()
	if err != nil {
		return status.Errorf(codes.Internal, "prepare update: %v", err)
	}

	_, err = stmt.ExecContext(ctx, password, u.Pb.GetId())
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
