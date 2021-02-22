package model

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"user-service/internal/pkg/app"
	"user-service/pb/users"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Group model
type Group struct {
	Pb         users.Group
	WithAccess bool
}

// Get func
func (u *Group) Get(ctx context.Context, db *sql.DB) error {
	query := `
		SELECT groups.id, groups.company_id, groups.name, groups.is_mutable, 
			json_agg(DISTINCT jsonb_build_object(
				'id', access.id,
				'name', access.name
			)) as access
		FROM groups 
		LEFT JOIN access_groups ON groups.id = access_groups.group_id
		LEFT JOIN access ON access_groups.access_id = access.id 
		WHERE groups.id = $1
		GROUP BY groups.id, groups.company_id, groups.name, groups.is_mutable`
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare statement: %v", err)
	}
	defer stmt.Close()

	var tmpAccess string
	err = stmt.QueryRowContext(ctx, u.Pb.GetId()).Scan(&u.Pb.Id, &u.Pb.CompanyId, &u.Pb.Name, &u.Pb.IsMutable, &tmpAccess)

	if err == sql.ErrNoRows {
		return status.Errorf(codes.NotFound, "Query Raw: %v", err)
	}

	if err != nil {
		return status.Errorf(codes.Internal, "Query Raw: %v", err)
	}

	err = json.Unmarshal([]byte(tmpAccess), &u.Pb.Access)
	if err != nil {
		return status.Errorf(codes.Internal, "unmarshal access: %v", err)
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

	u.Pb.CreatedAt = now.String()
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
	u.Pb.UpdatedAt = now.String()

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

// ListQuery builder
func (u *Group) ListQuery(ctx context.Context, db *sql.DB, in *users.ListGroupRequest) (string, []interface{}, *users.GroupPaginationResponse, error) {
	var paginationResponse users.GroupPaginationResponse
	query := `
		SELECT groups.id, groups.company_id, groups.name, groups.is_mutable, 
			json_agg(DISTINCT jsonb_build_object(
				'id', access.id,
				'name', access.name
			)) as access
		FROM groups 
		LEFT JOIN access_groups ON groups.id = access_groups.group_id
		LEFT JOIN access ON access_groups.access_id = access.id 
	`
	where := []string{"groups.company_id = $1"}
	paramQueries := []interface{}{ctx.Value(app.Ctx("companyID")).(string)}

	if len(in.GetPagination().GetSearch()) > 0 {
		paramQueries = append(paramQueries, in.GetPagination().GetSearch())
		where = append(where, fmt.Sprintf(`(groups.name ILIKE $%d)`, len(paramQueries)))
	}

	{
		qCount := `
			SELECT COUNT(*) 
			FROM groups 
			LEFT JOIN access_groups ON groups.id = access_groups.group_id
			LEFT JOIN access ON access_groups.access_id = access.id 
		`
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

	query += ` GROUP BY groups.id, groups.company_id, groups.name, groups.is_mutable`

	if len(in.GetPagination().GetOrderBy()) == 0 || !(in.GetPagination().GetOrderBy() == "groups.name") {
		if in.GetPagination() == nil {
			in.Pagination = &users.Pagination{OrderBy: "groups.created_at"}
		} else {
			in.GetPagination().OrderBy = "groups.created_at"
		}
	}

	query += ` ORDER BY ` + in.GetPagination().GetOrderBy() + ` ` + in.GetPagination().GetSort().String()

	if in.GetPagination().GetLimit() > 0 {
		query += fmt.Sprintf(` LIMIT $%d OFFSET $%d`, (len(paramQueries) + 1), (len(paramQueries) + 2))
		paramQueries = append(paramQueries, in.GetPagination().GetLimit(), in.GetPagination().GetOffset())
	}

	return query, paramQueries, &paginationResponse, nil
}
