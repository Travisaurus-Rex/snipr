package db

import (
	"context"
	"database/sql"
	"errors"
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
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user, password, host, port, name,
	)
	var err error
	Pool, err = pgxpool.New(context.TODO(), connString)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	fmt.Println("Connected to postgres")
}

func InsertURL(shortCode string, longURL string) (URL, error) {
	var u URL

	if Pool == nil {
		return u, errors.New("database connection not available")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := `
		INSERT INTO urls (short_code, long_url)
		VALUES ($1, $2)
		RETURNING id, short_code, long_url, created_at;
	`

	row := Pool.QueryRow(ctx, query, shortCode, longURL)
	err := row.Scan(&u.ID, &u.ShortCode, &u.LongURL, &u.CreatedAt)
	if err != nil {
		return u, err
	}

	return u, nil
}

func GetURLByCode(shortCode string) (URL, error) {
	var u URL

	if Pool == nil {
		return u, errors.New("database connection not available")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := `
		SELECT id, short_code, long_url, created_at
		FROM urls
		WHERE short_code = $1
	`

	row := Pool.QueryRow(ctx, query, shortCode)

	err := row.Scan(&u.ID, &u.ShortCode, &u.LongURL, &u.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return u, errors.New("no URL found for that short code")
		}
		return u, err
	}

	return u, nil
}
