package database

import (
	"Corap-web/models"
	"fmt"
	"log"
	"os"
	"time"

	"database/sql"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var (
	db *sql.DB
	// mu sync.Mutex
)

func Connect() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		dbUser, dbPassword, dbName, dbHost, dbPort)
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected with Database")
}

func GetDevices() []models.Device {
	rows, err := db.Query(`SELECT deveui, name, hashedname, is_defect FROM device d
							join analyse_device a on a.device_id = d.deveui
							order by name`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var devices []models.Device

	for rows.Next() {
		var device models.Device
		err := rows.Scan(&device.Deveui, &device.Name, &device.Hashedname, &device.IsDefect)
		if err != nil {
			log.Fatal(err)
		}
		devices = append(devices, device)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return devices
}

func GetDevice(deveui string) (models.Device, error) {
	row := db.QueryRow(`SELECT deveui, name, hashedname, is_defect FROM device d
							join analyse_device a on a.device_id = d.deveui
							WHERE deveui = $1`, deveui)
	var device models.Device
	err := row.Scan(&device.Deveui, &device.Name, &device.Hashedname, &device.IsDefect)
	if err != nil {
		if err == sql.ErrNoRows {
			return device, err
		}
		log.Fatal(err)
	}

	return device, nil
}

func GetDeviceScrapes(deveui string, plotType models.PlotType) ([]float32, []time.Time) {
	rows, err := db.Query(fmt.Sprintf(`SELECT %s, time_scraped FROM scrape
							WHERE deveui = $1
							AND time_scraped >= NOW() - '1 day'::INTERVAL
							ORDER BY time_scraped`, plotType), deveui)
	if err != nil {
		log.Fatal(err)
	}
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
	rows, err := db.Query("SELECT id, next_run_time, job_state FROM apscheduler_jobs")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var jobs []models.Job
	var float_time float64
	for rows.Next() {
		var job models.Job
		err := rows.Scan(&job.Id, &float_time, &job.JobState)
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

func Insert(user *models.User) {
}

func Get() *sql.DB {
	return db
}
