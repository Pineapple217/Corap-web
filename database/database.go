package database

import (
	"Corap-web/models"
	"context"
	"fmt"
	"os"
	"runtime"
	"time"

	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var (
	db     *pgxpool.Pool
	dbName string
	// mu sync.Mutex
)

func Connect() {
	err := godotenv.Load()
	if err != nil {
		log.Info("No .env file found.")
	}
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

func GetDevices() []models.Device {
	rows, err := db.Query(context.Background(), `WITH RankedScrapeData AS (
							SELECT d.deveui, d.name, d.hashedname, a.is_defect, s.temp, s.co2, s.humidity, s.time_scraped,
								ROW_NUMBER() OVER (PARTITION BY d.deveui ORDER BY s.time_scraped DESC) AS rn
							FROM device d
							JOIN scrape s ON d.deveui = s.deveui
							JOIN analyse_device a ON a.device_id = d.deveui
							)
							SELECT deveui, name, hashedname, is_defect, temp, co2, humidity
							FROM RankedScrapeData
							WHERE rn = 1;`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var devices []models.Device

	for rows.Next() {
		var device models.Device
		err := rows.Scan(&device.Deveui, &device.Name, &device.Hashedname, &device.IsDefect, &device.Temp, &device.Co2, &device.Humidity)
		if err != nil {
			log.Fatal(err)
		}
		devices = append(devices, device)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatal(err)
	}
	return devices
}

func GetDevice(deveui string) (models.Device, error) {
	row := db.QueryRow(context.Background(), `SELECT d.deveui, name, hashedname, is_defect, temp, co2, humidity FROM device d
						join analyse_device a on a.device_id = d.deveui
						join scrape s on d.deveui = s.deveui
						WHERE d.deveui = $1
						ORDER BY time_scraped DESC
						LIMIT 1;`, deveui)
	var device models.Device
	err := row.Scan(&device.Deveui, &device.Name, &device.Hashedname, &device.IsDefect,
		&device.Temp, &device.Co2, &device.Humidity)
	if err != nil {
		if err == pgx.ErrNoRows {
			return device, err
		}
		log.Fatal(err)
	}

	return device, nil
}

func GetDeviceScrapes(deveui string, plotType models.PlotType, dateRange int) ([]float32, []time.Time) {
	rows, err := db.Query(context.Background(), fmt.Sprintf(`SELECT %s, time_updated FROM scrape
							WHERE deveui = $1
							AND time_scraped >= NOW() - '%d day'::INTERVAL
							ORDER BY time_scraped`, plotType, dateRange), deveui)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var data float32
	var timestamp time.Time
	var datas []float32
	var timestamps []time.Time
	for rows.Next() {
		err := rows.Scan(&data, &timestamp)
		if err != nil {
			log.Fatal(err)
		}
		datas = append(datas, data)
		timestamps = append(timestamps, timestamp)

	}
	return datas, timestamps
}

func GetSchedulerJobs() []models.Job {
	rows, err := db.Query(context.Background(), "SELECT id, next_run_time FROM apscheduler_jobs")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var jobs []models.Job
	var float_time float64
	for rows.Next() {
		var job models.Job
		err := rows.Scan(&job.Id, &float_time)
		if err != nil {
			log.Fatal(err)
		}
		job.NextRunTime = time.Unix(int64(float_time), 0)
		jobs = append(jobs, job)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return jobs
}

func Get() *pgxpool.Pool {
	return db
}
