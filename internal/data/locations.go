package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/rytwalker/kagubird-api/internal/validator"
)

type Location struct {
	ID            int64     `json:"id"`
	Name          string    `json:"name"`
	Address       string    `json:"address"`
	Lat           float64   `json:"lat"`
	Lng           float64   `json:"lng"`
	GooglePlaceID string    `json:"google_place_id"`
	Website       string    `json:"website"`
	Phone         string    `json:"phone"`
	ActivityID    int64     `json:"activity"`
	Version       int32     `json:"version"`
	CreatedAt     time.Time `json:"-"`
	UpdatedAt     time.Time `json:"-"`
}

type LocationModel struct {
	DB *sql.DB
}

func (m LocationModel) Insert(location *Location) error {
	query := `
    INSERT INTO locations (name, address, lat, lng, google_place_id, website, phone, activity_id)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    RETURNING id, created_at, version`

	args := []any{location.Name, location.Address, location.Lat, location.Lng, location.GooglePlaceID, location.Website, location.Phone, location.ActivityID}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&location.ID, &location.CreatedAt, &location.Version)
}

func (m LocationModel) GetAllByActivity(activity_id int64) ([]*Location, error) {
	query := `
    SELECT  id, name, address, lat, lng, google_place_id, website, phone, version 
    FROM locations
    WHERE activity_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, activity_id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	locations := []*Location{}

	// rows.Next iterates
	for rows.Next() {
		var location Location

		err := rows.Scan(
			&location.ID,
			&location.Name,
			&location.Address,
			&location.Lat,
			&location.Lng,
			&location.GooglePlaceID,
			&location.Website,
			&location.Phone,
			&location.Version,
		)

		if err != nil {
			return nil, err
		}

		location.ActivityID = activity_id

		locations = append(locations, &location)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return locations, nil
}
func ValidateLocation(v *validator.Validator, location *Location) {
	// name validations
	v.Check(location.Name != "", "name", "must be provided")
	v.Check(len(location.Name) <= 500, "name", "must not be more than 500 bytes long")

	// address validations
	v.Check(location.Address != "", "address", "must be provided")
	v.Check(len(location.Address) <= 500, "city", "must not be more than 500 bytes long")

	// google_place_id validations
	v.Check(location.GooglePlaceID != "", "google_place_id", "must be provided")
	v.Check(len(location.GooglePlaceID) <= 500, "google_place_id", "must not be more than 500 bytes long")

	// lat validations
	v.Check(location.Lat != 0, "lat", "must be provided")

	// lng validations
	v.Check(location.Lng != 0, "lng", "must be provided")

	// activity_id validations
	v.Check(location.ActivityID != 0, "activity_id", "must be provided")

}
