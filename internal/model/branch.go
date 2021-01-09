package model

import (
	"context"
	"database/sql"
	users "user-service/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Branch model
type Branch struct {
	Pb users.Branch
}

// Get func
func (u *Branch) Get(ctx context.Context, db *sql.DB) error {
	query := `SELECT id, company_id, region_id, name, code, address, city, province, npwp, phone, pic, pic_phone 
	FROM branches WHERE id = $1`
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare statement: %v", err)
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, u.Pb.GetId()).Scan(
		&u.Pb.Id, &u.Pb.CompanyId, &u.Pb.RegionId, &u.Pb.Name, &u.Pb.Code, &u.Pb.Address, &u.Pb.City, &u.Pb.Province,
		&u.Pb.Npwp, &u.Pb.Phone, &u.Pb.Pic, &u.Pb.PicPhone,
	)

	if err == sql.ErrNoRows {
		return status.Errorf(codes.NotFound, "Query Raw: %v", err)
	}

	if err != nil {
		return status.Errorf(codes.Internal, "Query Raw: %v", err)
	}

	return nil
}
