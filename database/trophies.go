package database

import (
	"Corap-web/models"
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2/log"
)

func GetTrophies(tType models.DataType) [3]models.Trophy {
	rows, err := db.Query(context.Background(), fmt.Sprintf(`
		SELECT sub.deveui, d.name, s.time_updated, s.temp, s.co2, s.humidity
		FROM (
			SELECT deveui, MAX(%s) AS max_d
			FROM scrape
			GROUP BY deveui
		) AS sub
		JOIN scrape s ON sub.deveui = s.deveui AND sub.max_d = s.%s
		JOIN device d on d.deveui = sub.deveui
		ORDER BY sub.max_d DESC
		limit 3;
	`, tType, tType))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var trophies [3]models.Trophy
	i := 0
	for rows.Next() {
		var trophy models.Trophy
		err := rows.Scan(&trophy.Deveui, &trophy.DevName, &trophy.TimeUpdated,
			&trophy.Temp, &trophy.Co2, &trophy.Humidity)
		if err != nil {
			log.Fatal(err)
		}
		trophies[i] = trophy
		i++
	}
	return trophies
}
