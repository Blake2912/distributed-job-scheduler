package main

import (
	"fmt"
	"log"
	"os"

	"example.com/sql_dal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// TODO:: Ensure LoadEnv is called only in dev envs.. need to create separate envs for dev and prod
func main() {
	config.LoadEnv()
	fmt.Println("Starting to connect to database")

	db, err := Connect()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to Postgres via GORM")
	
	fmt.Println(db.DB())

}

func Connect() (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSLMODE"),
	)

	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
