package main

import (
	"flag"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"

	"user-service/internal/config"
	"user-service/internal/pkg/db/postgres"
	"user-service/internal/pkg/log/logruslog"
	"user-service/internal/schema"
)

func main() {
	if _, ok := os.LookupEnv("APP_ENV"); !ok {
		_, err := os.Stat(".env.prod")
		if os.IsNotExist(err) {
			config.Setup(".env")
		} else {
			config.Setup(".env.prod")
		}
	}

	// =========================================================================
	// Logging
	log := logruslog.Init()
	if err := run(log); err != nil {
		log.Printf("error: shutting down: %s", err)
		os.Exit(1)
	}
}

func run(log *logrus.Entry) error {
	// =========================================================================
	// App Starting

	log.Printf("main : Started")
	defer log.Println("main : Completed")

	// =========================================================================

	// Start Database

	db, err := postgres.Open()
	if err != nil {
		return fmt.Errorf("connecting to db: %v", err)
	}
	defer db.Close()

	// Handle cli command
	flag.Parse()

	switch flag.Arg(0) {
	case "migrate":
		if err := schema.Migrate(db); err != nil {
			return fmt.Errorf("applying migrations: %v", err)
		}
		log.Println("Migrations complete")
		return nil

	case "seed":
		if err := schema.Seed(db); err != nil {
			return fmt.Errorf("seeding database: %v", err)
		}
		log.Println("Seed data complete")
		return nil
	}

	return nil
}
