package schema

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func initSeed(ctx context.Context, tx *sql.Tx) error {
	packageFeatureID := uuid.New().String()
	companyID := uuid.New().String()
	userID := uuid.New().String()
	groupID := uuid.New().String()
	accessID := uuid.New().String()

	// seed package features
	stmt, err := tx.PrepareContext(ctx,
		`INSERT INTO package_features (id, name) VALUES ($1, 'ALL'), ($2, 'SIMPLE'), ($3, 'CUSTOME')`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, packageFeatureID, uuid.New().String(), uuid.New().String())
	if err != nil {
		return err
	}

	// seed company
	query := `
		INSERT INTO companies (id, name, code, address, city, province, phone, pic, pic_phone, package_of_feature_id, updated_by)
		VALUES ($1, 'Wiradata Sistem', 'WIRA', 'Pondok Aren', 'Tangerang Selatan', 'Banten', '08122222222', 'Jacky', '08133333333', $2, $3)
	`

	stmt, err = tx.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, companyID, packageFeatureID, userID)
	if err != nil {
		return err
	}

	// seed group
	query = `
		INSERT INTO groups (id, company_id, is_mutable, name, created_by, updated_by)
		VALUES ($1, $2, false, 'Super Admin', $3, $3)
	`

	stmt, err = tx.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, groupID, companyID, userID)
	if err != nil {
		return err
	}

	// seed user
	password, err := bcrypt.GenerateFromPassword([]byte("1234"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	query = `
		INSERT INTO users (id, company_id, group_id, username, name, email, password)
		VALUES ($1, $2, $3, 'wira-admin', 'Administrator', 'rijal.asep.nugroho@gmail.com', $4)
	`

	stmt, err = tx.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, userID, companyID, groupID, password)
	if err != nil {
		return err
	}

	// seed access
	query = `
		INSERT INTO access (id, name, created_by, updated_by)
		VALUES ($1, 'root', $2, $2)
	`

	stmt, err = tx.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, accessID, userID)
	if err != nil {
		return err
	}

	// seed grant access
	query = `
		INSERT INTO access_groups (id, group_id, access_id, created_by, updated_by)
		VALUES ($1, $2, $3, $4, $4)
	`

	stmt, err = tx.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, uuid.New().String(), groupID, accessID, userID)
	if err != nil {
		return err
	}

	return nil
}

// seeds is a string constant containing all of the queries needed to get the
// db seeded to a useful state for development.
//
// Using a constant in a .go file is an easy way to ensure the queries are part
// of the compiled executable and avoids pathing issues with the working
// directory. It has the downside that it lacks syntax highlighting and may be
// harder to read for some cases compared to using .sql files. You may also
// consider a combined approach using a tool like packr or go-bindata.
//
// Note that database servers besides PostgreSQL may not support running
// multiple queries as part of the same execution so this single large constant
// may need to be broken up.

// Seed runs the set of seed-data queries against db. The queries are ran in a
// transaction and rolled back if any fail.
func Seed(db *sql.DB) error {
	seeds := []string{}
	ctx := context.Background()

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	err = initSeed(ctx, tx)
	if err != nil {
		tx.Rollback()
		fmt.Println("error execute initSeed")
		return err
	}

	for _, seed := range seeds {
		_, err = tx.ExecContext(ctx, seed)
		if err != nil {
			tx.Rollback()
			fmt.Println("error execute seed")
			return err
		}
	}

	return tx.Commit()
}
