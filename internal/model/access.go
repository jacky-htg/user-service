package model

import (
	"context"
	"database/sql"
	"encoding/json"
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
		//pbAccess.ParentId = parent

		list = append(list, &pbAccess)
	}

	return list, nil
}

func (u *Access) List(ctx context.Context, db *sql.DB) (*users.ListAccessResponse, error) {
	var output users.ListAccessResponse
	rows, err := db.QueryContext(ctx, `SELECT id, name FROM access WHERE level = 1`)

	if err == sql.ErrNoRows {
		return &output, status.Errorf(codes.NotFound, "Query Raw: %v", err)
	}

	if err != nil {
		return &output, status.Errorf(codes.Internal, "Query Raw: %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var pbAccess users.AccessLevel1
		err = rows.Scan(&pbAccess.Id, &pbAccess.Name)
		if err != nil {
			return &output, status.Errorf(codes.Internal, "Scan: %v", err)
		}

		pbAccess.Children, err = u.ListChildren(ctx, db, pbAccess.Id)
		if err != nil {
			return nil, err
		}
		output.Access = append(output.Access, &pbAccess)
	}

	if rows.Err() != nil {
		return nil, status.Error(codes.Internal, rows.Err().Error())
	}

	return &output, nil
}

func (u *Access) ListChildren(ctx context.Context, db *sql.DB, parentId string) ([]*users.AccessLevel2, error) {
	var err error
	var output []*users.AccessLevel2

	query := `
	SELECT access2.id, access2.name,
	JSON_AGG(DISTINCT JSONB_BUILD_OBJECT(
		'id', access3.id,
		'name', access3.name
	)) as json
	FROM access access2
	LEFT JOIN access access3 ON access3.level = 3 and access3.parent_id = access2.id
	WHERE access2.parent_id  = $1 AND access2.level = 2
	GROUP BY access2.id
	`
	rows, err := db.QueryContext(ctx, query, parentId)
	if err != nil {
		return output, status.Error(codes.Internal, err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var pbAccess users.AccessLevel2
		var tempJson string
		err = rows.Scan(&pbAccess.Id, &pbAccess.Name, &tempJson)
		if err != nil {
			return output, status.Errorf(codes.Internal, "scan data: %v", err)
		}
		err = json.Unmarshal([]byte(tempJson), &pbAccess.Children)
		if err != nil {
			return output, status.Errorf(codes.Internal, "unmarshal access: %v", err)
		}

		output = append(output, &pbAccess)
	}

	if rows.Err() != nil {
		return output, status.Error(codes.Internal, rows.Err().Error())
	}

	return output, nil
}
