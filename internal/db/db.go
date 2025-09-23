package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var Pool *pgxpool.Pool

func Connect() {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")

	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		user, password, host, port, name,
	)
	var err error
	Pool, err = pgxpool.New(context.TODO(), connString)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	fmt.Println("Connected to postgres")
}
