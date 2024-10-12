package database

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/Pineapple217/Corap-web/pkg/helper"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/joho/godotenv/autoload"
)

type Database struct {
	pool *pgxpool.Pool
}

func NewDatabase() *Database {
	ctx := context.Background()

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_DATABASE")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	slog.Info("Starting database", "host", dbHost, "db", dbName)
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s",
		dbUser, dbPassword, dbName, dbHost, dbPort)
	pool, err := pgxpool.New(ctx, connStr)
	helper.MaybeDie(err, "Database err")

	err = pool.Ping(ctx)
	helper.MaybeDie(err, "Failed to connect to database")

	db := &Database{
		pool: pool,
	}

	return db
}
