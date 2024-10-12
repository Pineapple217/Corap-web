package database

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type Device struct {
	Deveui     string  `json:"deveui"`
	Name       string  `json:"name"`
	Hashedname string  `json:"hashedname"`
	Temp       float32 `json:"temp"`
	Co2        int     `json:"co2"`
	Humi       int     `json:"humi"`
}

type DeviceAnalysis struct {
	Deveui     string     `json:"deveui"`
	Name       string     `json:"name"`
	Hashedname string     `json:"hashedname"`
	Temp       float32    `json:"temp"`
	Co2        int        `json:"co2"`
	Humi       int        `json:"humi"`
	IsDefect   bool       `json:"isDefect"`
	Analysis   []Analysis `json:"analysis"`
}

type Analysis struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type DataType string

const (
	Temp DataType = "temp"
	Co2  DataType = "co2"
	Humi DataType = "humi"
)

func (d *Device) Url() string {
	return "https://education.thingsflow.eu/IAQ/DeviceByQR?hashedname=" + d.Hashedname
}

func (d *DeviceAnalysis) Url() string {
	return "https://education.thingsflow.eu/IAQ/DeviceByQR?hashedname=" + d.Hashedname
}

func (db *Database) AllDevices(ctx context.Context) ([]*Device, error) {
	rows, err := db.pool.Query(ctx, `
		WITH LatestScrape AS (
			SELECT
				d.deveui,
				MAX(s.scraped_at) AS latest_time_scraped
			FROM
				devices d
				LEFT JOIN scrapes s ON d.deveui = s.deveui AND s.scraped_at >= NOW() - '1 hour'::INTERVAL
			GROUP BY
				d.deveui
		),
		RankedScrapeData AS (
			SELECT
				d.deveui,
				d.name,
				d.hashedname,
				COALESCE(s.temp, -1) AS temp,
				COALESCE(s.co2, -1) AS co2,
				COALESCE(s.humi, -1) AS humidity,
				COALESCE(s.scraped_at, NOW()) AS scraped_at,
				ROW_NUMBER() OVER (PARTITION BY d.deveui ORDER BY s.scraped_at DESC) AS rn
			FROM
				devices d
				LEFT JOIN scrapes s ON d.deveui = s.deveui AND s.scraped_at = (SELECT latest_time_scraped FROM LatestScrape WHERE deveui = d.deveui)
		)
		SELECT
			deveui,
			name,
			hashedname,
			temp,
			co2,
			humidity
		FROM
			RankedScrapeData
		WHERE
			rn = 1;
	`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var devices []*Device

	for rows.Next() {
		var device Device
		err := rows.Scan(&device.Deveui, &device.Name, &device.Hashedname, &device.Temp, &device.Co2, &device.Humi)
		if err != nil {
			return nil, err
		}
		devices = append(devices, &device)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return devices, nil
}

func (db *Database) AllDevicesAnalysis(ctx context.Context) ([]*DeviceAnalysis, error) {
	rows, err := db.pool.Query(ctx, `
		SELECT * FROM public.device_analysis_summary;
	`)
	if err != nil {
		return nil, err
	}
	rowNames := rows.FieldDescriptions()

	defer rows.Close()

	var devices []*DeviceAnalysis

	for rows.Next() {
		var device DeviceAnalysis
		v, err := rows.Values()
		if err != nil {
			return nil, err
		}
		device.Deveui = v[0].(string)
		device.Name = v[1].(string)
		device.Hashedname = v[2].(string)
		f, err := v[3].(pgtype.Numeric).Float64Value()
		if err != nil {
			return nil, err
		}
		device.Temp = float32(f.Float64)
		device.Co2 = int(v[4].(int32))
		device.Humi = int(v[5].(int32))
		b, err := strconv.ParseBool(v[6].(string))
		if err != nil {
			return nil, err
		}
		device.IsDefect = b

		offset := 7
		device.Analysis = []Analysis{}
		for i, a := range v[offset:] {
			device.Analysis = append(device.Analysis, Analysis{
				Value: a.(string),
				Name:  rowNames[offset+i].Name,
			})
		}

		devices = append(devices, &device)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return devices, nil
}

func (db *Database) DeviceById(ctx context.Context, id string) (*Device, error) {
	row := db.pool.QueryRow(ctx,
		`SELECT d.deveui, name, hashedname, temp, co2, humi FROM devices d
		join scrapes s on d.deveui = s.deveui
		WHERE d.deveui = $1
		ORDER BY s.scraped_at DESC
		LIMIT 1;`, id)
	var device Device
	err := row.Scan(&device.Deveui, &device.Name, &device.Hashedname,
		&device.Temp, &device.Co2, &device.Humi)
	if err != nil {
		return nil, err
	}

	return &device, nil
}

func (db *Database) DeviceHistory(ctx context.Context, deveui string,
	dateType DataType, dayRange int) ([][]any, error) {

	rows, err := db.pool.Query(ctx,
		fmt.Sprintf(`SELECT %s, scraped_at FROM scrapes
		WHERE deveui = $1
		AND scraped_at >= NOW() - '%d day'::INTERVAL
		ORDER BY scraped_at`, dateType, dayRange), deveui)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var data any
	var timestamp time.Time
	var datas [][]any
	for rows.Next() {
		err := rows.Scan(&data, &timestamp)
		if err != nil {
			return nil, err
		}
		datas = append(datas, []any{timestamp, data})

	}
	return datas, nil
}

func FormatDataType(plotTypeStr string) (DataType, error) {
	switch plotTypeStr {
	case string(Temp):
		return Temp, nil
	case string(Co2):
		return Co2, nil
	case string(Humi):
		return Humi, nil
	default:
		return "", errors.New("string is not a valid datatype")
	}
}
