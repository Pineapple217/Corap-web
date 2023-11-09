package database

import (
	"context"
	"fmt"
	"os"
	"runtime"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	db     *pgxpool.Pool
	dbName string
)

func Connect() {
	var err error
	maxConn := runtime.NumCPU() * 4
	if fiber.IsChild() {
		maxConn = 5
	}
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName = os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable pool_max_conns=%d",
		dbUser, dbPassword, dbName, dbHost, dbPort, maxConn)
	db, err = pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Connected with Database")
}

func Get() *pgxpool.Pool {
	return db
}
