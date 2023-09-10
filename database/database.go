package database

import (
	"Corap-web/models"
	"fmt"
	"log"
	"time"

	"database/sql"

	_ "github.com/lib/pq"
)

var (
	db *sql.DB
	// mu sync.Mutex
)

// Connect with database
func Connect() {
	// connStr := "user=read_only password=password dbname=corap host=localhost port=9979 sslmode=disable"
	connStr := "user=corap password=root dbname=corap host=localhost port=5432 sslmode=disable"
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	db.Ping()
	fmt.Println("Connected with Database")
}

func GetDivices() []models.Device {
	// rows, err := db.Query("SELECT deveui, name, hashedname FROM device")
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
