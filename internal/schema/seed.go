package schema

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func initSeed(ctx context.Context, tx *sql.Tx) error {
	// seed package features
	query := fmt.Sprintf(`
		INSERT INTO package_features (id, name) VALUES
		(%s, 'ALL'),
		(%s, 'SIMPLE'),
		(%s, 'CUSTOME')
	`, uuid.New().String(), uuid.New().String(), uuid.New().String())

	_, err := tx.ExecContext(ctx, query)
	if err != nil {
		return err
	}

	// seed company
	var packageFeatureID, companyID, userID string
	err = tx.QueryRowContext(ctx, `SELECT id FROM package_features WHERE name='ALL'`).Scan(&packageFeatureID)
	if err != nil {
		return err
	}

	query = fmt.Sprintf(`
		INSERT INTO companies (id, name, code, address, city, province, phone, pic, pic_phone, package_of_feature_id)
		VALUES (%s, 'Wiradata Sistem', 'WIRA', 'Pondok Aren', 'Tangerang Selatan', 'Banten', '08122222222', 'Jacky', '08133333333', %s)
		RETURNING id
	`, uuid.New().String(), packageFeatureID)

	err = tx.QueryRowContext(ctx, query).Scan(&companyID)
	if err != nil {
		return err
	}

	// seed user
	password, err := bcrypt.GenerateFromPassword([]byte("1234"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	query = fmt.Sprintf(`
		INSERT INTO users (id, company_id, username, name, email, password)
		VALUES (%s, %s, 'wira-admin', 'Administrator', 'rijal.asep.nugroho@gmail.com', %s)
		RETURNING id
	`, uuid.New().String(), companyID, password)

	err = tx.QueryRowContext(ctx, query).Scan(&userID)
	if err != nil {
		return err
	}

	// update company
	query = fmt.Sprintf(`UPDATE companies SET update_by = %s WHERE id=%s`, userID, companyID)
	_, err = tx.ExecContext(ctx, query)
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
