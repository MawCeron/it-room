package repo

import (
	"database/sql"

	"github.com/MawCeron/it-room/internal/models"
)

type LocationRepo struct{ db *sql.DB }

func NewLocationRepo(db *sql.DB) *LocationRepo {
	return &LocationRepo{db: db}
}

func (r *LocationRepo) List() ([]*models.Location, error) {
	rows, err := r.db.Query(`SELECT location_id, name, "type"
FROM locations;`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*models.Location
	for rows.Next() {
		var l models.Location

		if err := rows.Scan(&l.LocationID, &l.Name, &l.Type); err != nil {
			return nil, err
		}
		out = append(out, &l)
	}

	return out, nil
}
