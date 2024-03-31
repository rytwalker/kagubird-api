package data

import (
	"database/sql"
	"errors"
	"time"

	"github.com/rytwalker/kagubird-api/internal/validator"
)

type Trip struct {
	ID            int64     `json:"id"`
	CreatedAt     time.Time `json:"-"`
	Name          string    `json:"name"`
	City          string    `json:"city"`
	StateCode     string    `json:"state_code"`
	GooglePlaceID string    `json:"google_place_id"`
	Lat           float64   `json:"lat"`
	Lng           float64   `json:"lng"`
	StartDate     time.Time `json:"start_date"`
	EndDate       time.Time `json:"end_date"`
	Version       int32     `json:"version"`
}

type TripModel struct {
	DB *sql.DB
}

func (t TripModel) Insert(trip *Trip) error {
	query := `
    INSERT INTO trips (name, city, state_code, google_place_id, lat, lng, start_date, end_date)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    RETURNING id, created_at, version`

	args := []any{trip.Name, trip.City, trip.StateCode, trip.GooglePlaceID, trip.Lat, trip.Lng, trip.StartDate.UTC(), trip.EndDate.UTC()}

	return t.DB.QueryRow(query, args...).Scan(&trip.ID, &trip.CreatedAt, &trip.Version)
}

func (t TripModel) Get(id int64) (*Trip, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
    SELECT id, created_at, name, city, state_code, google_place_id, lat, lng, start_date, end_date
    FROM trips
    WHERE id = $1`

	var trip Trip

	err := t.DB.QueryRow(query, id).Scan(
		&trip.ID,
		&trip.CreatedAt,
		&trip.Name,
		&trip.City,
		&trip.StateCode,
		&trip.GooglePlaceID,
		&trip.Lat,
		&trip.Lng,
		&trip.StartDate,
		&trip.EndDate,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &trip, nil
}

func (t TripModel) Update(trip *Trip) error {
	query := `
    UPDATE trips
    SET name = $1, city = $2, state_code = $3, google_place_id = $4, lat = $5, lng = $6, start_date = $7, end_date = $8, version = version + 1
    WHERE id = $9
    RETURNING version`

	args := []any{
		trip.Name,
		trip.City,
		trip.StateCode,
		trip.GooglePlaceID,
		trip.Lat,
		trip.Lng,
		trip.StartDate,
		trip.EndDate,
		trip.ID,
	}

	return t.DB.QueryRow(query, args...).Scan(&trip.Version)
}

func (t TripModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
    DELETE FROM trips
    WHERE id = $1`

	result, err := t.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func ValidateTrip(v *validator.Validator, trip *Trip) {
	// name validations
	v.Check(trip.Name != "", "name", "must be provided")
	v.Check(len(trip.Name) <= 500, "name", "must not be more than 500 bytes long")

	// city validations
	v.Check(trip.City != "", "city", "must be provided")
	v.Check(len(trip.City) <= 500, "city", "must not be more than 500 bytes long")

	// state_code validations
	v.Check(trip.StateCode != "", "state_code", "must be provided")
	v.Check(len(trip.StateCode) <= 500, "state_code", "must not be more than 500 bytes long")

	// google_place_id validations
	v.Check(trip.GooglePlaceID != "", "google_place_id", "must be provided")
	v.Check(len(trip.GooglePlaceID) <= 500, "google_place_id", "must not be more than 500 bytes long")

	// lat validations
	v.Check(trip.Lat != 0, "lat", "must be provided")

	// lng validations
	v.Check(trip.Lng != 0, "lng", "must be provided")

	// start_date validations
	v.Check(!trip.StartDate.IsZero(), "start_date", "must be provided")
	v.Check(!trip.StartDate.Before(time.Now()), "start_date", "must be in the future")
	v.Check(trip.StartDate.Before(trip.EndDate), "start_date", "must be before end date")

	// end_date validations
	v.Check(!trip.EndDate.IsZero(), "end_date", "must be provided")
	v.Check(!trip.EndDate.Before(time.Now()), "end_date", "must be in the future")
	v.Check(trip.EndDate.After(trip.StartDate), "end_date", "must be after start date")
}
