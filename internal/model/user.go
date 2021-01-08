package model

import (
	"context"
	"database/sql"
	"encoding/json"
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

	if err != nil {
		return status.Errorf(codes.Internal, "Query Raw: %v", err)
	}

	u.Pb.RegionId = regionID.String
	u.Pb.BranchId = branchID.String

	err = bcrypt.CompareHashAndPassword([]byte(strPassword), []byte(password))
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "Invalid Password: %v", err)
	}

	err = json.Unmarshal([]byte(tmpAccess), &group.Access)
	if err != nil {
		return status.Errorf(codes.Internal, "unmarshal access: %v", err)
	}

	u.Pb.Group = &group

	return nil
}
