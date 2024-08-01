package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/jacky-htg/erp-pkg/db/postgres"
	"github.com/jacky-htg/erp-pkg/db/redis"
	"github.com/jacky-htg/user-service/internal/config"
	"github.com/jacky-htg/user-service/internal/route"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

const defaultPort = "8000"

func main() {
	// lookup and setup env
	if _, ok := os.LookupEnv("PORT"); !ok {
		config.Setup(".env")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// init log
	log := log.New(os.Stdout, "ERROR : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	// create postgres database connection
	db, err := postgres.Open()
	if err != nil {
		log.Fatalf("connecting to db: %v", err)
		return
	}

	defer db.Close()

	// create redis cache connection
	cache, err := redis.NewCache(context.Background(), os.Getenv("REDIS_ADDRESS"), os.Getenv("REDIS_PASSWORD"), 24*time.Hour)
	if err != nil {
		log.Fatalf("cannot create redis connection: %v", err)
		return
	}

	// listen tcp port
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return
	}

	grpcServer := grpc.NewServer()

	// routing grpc services
	route.GrpcRoute(grpcServer, db, log, cache)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
		return
	}
	fmt.Println("serve grpc on port: " + port)
}
