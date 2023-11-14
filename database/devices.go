package database

import (
	"Corap-web/models"
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/jackc/pgx/v5"
)

func GetDevices() []models.Device {
	rows, err := db.Query(context.Background(), `
		WITH LatestScrape AS (
			SELECT
				d.deveui,
				MAX(s.time_scraped) AS latest_time_scraped
			FROM
				device d
				LEFT JOIN scrape s ON d.deveui = s.deveui AND s.time_scraped >= NOW() - '1 hour'::INTERVAL
			GROUP BY
				d.deveui
		),
		RankedScrapeData AS (
			SELECT
				d.deveui,
				d.name,
				d.hashedname,
				a.is_defect,
				COALESCE(s.temp, -1) AS temp,
				COALESCE(s.co2, -1) AS co2,
				COALESCE(s.humidity, -1) AS humidity,
				COALESCE(s.time_scraped, NOW()) AS time_scraped,
				ROW_NUMBER() OVER (PARTITION BY d.deveui ORDER BY s.time_scraped DESC) AS rn
			FROM
				device d
				JOIN analysis_devices a ON a.device_id = d.deveui
				LEFT JOIN scrape s ON d.deveui = s.deveui AND s.time_scraped = (SELECT latest_time_scraped FROM LatestScrape WHERE deveui = d.deveui)
		)
		SELECT
			deveui,
			name,
			hashedname,
			is_defect,
			temp,
			co2,
			humidity
		FROM
			RankedScrapeData
		WHERE
			rn = 1;
	`)
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
						join analysis_devices a on a.device_id = d.deveui
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

func GetDeviceScrapes(deveui string, plotType models.DataType, dateRange int) ([]float32, []time.Time) {
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
