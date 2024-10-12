package database

import "context"

type Stat struct {
	Name  string
	Value string
}

func (db *Database) AllStats(ctx context.Context) ([]*Stat, error) {
	rows, err := db.pool.Query(ctx, `
		SELECT name, value
		FROM scrape_stats
		ORDER BY priority DESC;
	`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var stats []*Stat

	for rows.Next() {
		var stat Stat
		err := rows.Scan(&stat.Name, &stat.Value)
		if err != nil {
			return nil, err
		}
		stats = append(stats, &stat)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return stats, nil
}
