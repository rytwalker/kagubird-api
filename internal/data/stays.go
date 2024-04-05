package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/rytwalker/kagubird-api/internal/validator"
)

type Stay struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Address   string    `json:"address"`
	Lat       float64   `json:"lat"`
	Lng       float64   `json:"lng"`
	Link      string    `json:"link"`
	Phone     string    `json:"phone"`
	Type      string    `json:"type"`
	TripID    int64     `json:"trip"`
	Version   int32     `json:"version"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type StayModel struct {
	DB *sql.DB
}

func (m StayModel) Insert(stay *Stay) error {
	query := `
    INSERT INTO stays (name, address, start_time, end_time, lat, lng, link, phone, type, trip_id)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
    RETURNING id, created_at, version`

	args := []any{stay.Name, stay.Address, stay.StartTime, stay.EndTime, stay.Lat, stay.Lng, stay.Link, stay.Phone, stay.Type, stay.TripID}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&stay.ID, &stay.CreatedAt, &stay.Version)
}

func (m StayModel) GetAllByTrip(trip_id int64) ([]*Stay, error) {
	query := `
    SELECT  id, name, address, lat, lng,  start_time, end_time, link, phone, type, version 
    FROM stays
    WHERE activity_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, trip_id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	stays := []*Stay{}

	// rows.Next iterates
	for rows.Next() {
		var stay Stay

		err := rows.Scan(
			&stay.ID,
			&stay.Name,
			&stay.Address,
			&stay.Lat,
			&stay.Lng,
			&stay.StartTime,
			&stay.EndTime,
			&stay.Link,
			&stay.Phone,
			&stay.Type,
			&stay.Version,
		)

		if err != nil {
			return nil, err
		}

		stay.TripID = trip_id

		stays = append(stays, &stay)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return stays, nil
}
func ValidateStay(v *validator.Validator, stay *Stay) {
	// name validations
	v.Check(stay.Name != "", "name", "must be provided")
	v.Check(len(stay.Name) <= 500, "name", "must not be more than 500 bytes long")

	// address validations
	v.Check(stay.Address != "", "address", "must be provided")
	v.Check(len(stay.Address) <= 500, "address", "must not be more than 500 bytes long")
	v.Check(stay.Lat != 0, "lat", "must be provided")

	// lng validations
	v.Check(stay.Lng != 0, "lng", "must be provided")

	// activity_id validations
	v.Check(stay.TripID != 0, "trip_id", "must be provided")

}
