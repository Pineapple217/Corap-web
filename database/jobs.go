package database

import (
	"Corap-web/models"
	"context"
	"time"

	"github.com/gofiber/fiber/v2/log"
)

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
