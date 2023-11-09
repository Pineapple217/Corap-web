package database

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2/log"
)

func GetDatabaseSize() string {
	row := db.QueryRow(context.Background(), fmt.Sprintf("SELECT pg_size_pretty(pg_database_size('%s'));", dbName))
	var size string
	err := row.Scan(&size)
	if err != nil {
		log.Fatal(err)
	}
	return size
}

func GetScrapeCount() int {
	row := db.QueryRow(context.Background(), "SELECT COUNT(id) FROM scrape")
	var countStr string
	err := row.Scan(&countStr)
	if err != nil {
		log.Fatal(err)
	}
	count, err := strconv.ParseInt(countStr, 10, 0)
	if err != nil {
		log.Fatal(err)
	}
	return int(count)
}

func GetBatchCount() int {
	row := db.QueryRow(context.Background(), "SELECT MAX(batch_id) from scrape;")
	var batchCountStr string
	err := row.Scan(&batchCountStr)
	if err != nil {
		log.Fatal(err)
	}
	count, err := strconv.ParseInt(batchCountStr, 10, 0)
	if err != nil {
		log.Fatal(err)
	}
	return int(count)
}

func GetTimeLastScrape() time.Time {
	row := db.QueryRow(context.Background(), "SELECT MAX(time_scraped) from scrape;")
	var scrapeTime time.Time
	err := row.Scan(&scrapeTime)
	if err != nil {
		log.Fatal(err)
	}
	return scrapeTime
}
