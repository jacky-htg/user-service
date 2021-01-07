package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strconv"
)

// Open database commection
func Open() (*sql.DB, error) {
	var db *sql.DB
	port, err := strconv.Atoi(os.Getenv("POSTGRES_PORT"))
	if err != nil {
		return db, err
	}

	return sql.Open("postgres",
		fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			os.Getenv("POSTGRES_HOST"), port, os.Getenv("POSTGRES_USER"),
			os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"),
		),
	)
}

// StatusCheck returns nil if it can successfully talk to the database. It
// returns a non-nil error otherwise.
func StatusCheck(ctx context.Context, db *sql.DB) error {

	// Run a simple query to determine connectivity. The db has a "Ping" method
	// but it can false-positive when it was previously able to talk to the
	// database but the database has since gone away. Running this query forces a
	// round trip to the database.
	const q = `SELECT true`
	var tmp bool
	return db.QueryRowContext(ctx, q).Scan(&tmp)
}
