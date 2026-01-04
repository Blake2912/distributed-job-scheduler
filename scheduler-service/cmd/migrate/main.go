package main

import (
	"context"
	"log"

	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/sql_dal/config"
	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/sql_dal/migration"
)

func main() {
	log.Println("starting migration runner")

	ctx := context.Background()

	if err := config.ConnectDB(ctx); err != nil {
		log.Fatalf("unable to connect to database: %v", err)
	}
	defer config.CloseSqlConnection()

	if err := migration.Migrate(ctx, config.DB); err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	log.Println("migrations completed successfully")
}
