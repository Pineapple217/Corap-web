package models

import "time"

type Trophy struct {
	Deveui      string    `json:"deveui"`
	DevName     string    `json:"name"`
	Temp        float32   `json:"temp"`
	Co2         int       `json:"co2"`
	Humidity    int       `json:"humidity"`
	TimeUpdated time.Time `json:"time_updated"`
	// Hashedname string  `json:"hashedname"`
	// IsDefect   bool    `json:"is_defect"`
}
