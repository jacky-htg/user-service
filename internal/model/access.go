package model

import (
	"context"
	"database/sql"
	"user-service/pb/users"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Access struct
type Access struct {
	Pb users.Access
}

// Get Access
func (u *Access) Get(ctx context.Context, db *sql.DB) error {
	err := db.QueryRowContext(ctx, `SELECT id, name FROM access WHERE id = $1`, u.Pb.GetId()).Scan(&u.Pb.Id, &u.Pb.Name)

	if err == sql.ErrNoRows {
		return status.Errorf(codes.NotFound, "Query Raw: %v", err)
	}

	if err != nil {
		return status.Errorf(codes.Internal, "Query Raw: %v", err)
	}

	return nil
}

// GetRoot Access
func (u *Access) GetRoot(ctx context.Context, tx *sql.Tx, withChildren bool) error {
	err := tx.QueryRowContext(ctx, `SELECT id, name FROM access WHERE name = 'root'`).Scan(&u.Pb.Id, &u.Pb.Name)

	if err == sql.ErrNoRows {
		return status.Errorf(codes.NotFound, "Query Raw: %v", err)
	}

	if err != nil {
		return status.Errorf(codes.Internal, "Query Raw: %v", err)
	}

	/*if withChildren {
		u.Pb.Childrens, err = u.GetByParent(ctx, tx, u.Pb.GetId())
		if err != nil {
			return err
		}
	}*/

	return nil
}

// GetByParent Access
func (u *Access) GetByParent(ctx context.Context, tx *sql.Tx, parent string) ([]*users.Access, error) {
	var list []*users.Access

	rows, err := tx.QueryContext(ctx, `SELECT access.id, access.name FROM access WHERE parent_id = $1`, parent)
	if err != nil {
		return list, status.Error(codes.Internal, err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var pbAccess users.Access
		err = rows.Scan(&pbAccess.Id, &pbAccess.Name)
		if err != nil {
			return list, status.Error(codes.Internal, err.Error())
		}

		/*pbAccess.Childrens, err = u.GetByParent(ctx, tx, pbAccess.GetId())
		if err != nil {
			return list, err
		}*/
		pbAccess.ParentId = parent

		list = append(list, &pbAccess)
	}

	return list, nil
}
