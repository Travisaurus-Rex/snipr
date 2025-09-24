package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var Pool *pgxpool.Pool

type URL struct {
	ID        int
	ShortCode string
	LongURL   string
	CreatedAt time.Time
}

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

func InsertURL(shortCode string, longURL string) (*URL, error) {
	query := `
		INSERT INTO urls (short_code, long_url)
		VALUES ($1, $2)
		RETURNING id, short_code, long_url, created_at;
	`

	row := Pool.QueryRow(context.TODO(), query, shortCode, longURL)

	var u URL
	err := row.Scan(&u.ID, &u.ShortCode, &u.LongURL, &u.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &u, nil
}
