package models

import "time"

type Job struct {
	Id          string    `json:"id"`
	NextRunTime time.Time `json:"next_run_time"`
	JobState    string    `json:"job_state"`
}
