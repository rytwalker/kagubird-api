package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/rytwalker/kagubird-api/internal/validator"
)

type Activity struct {
	ID        int64       `json:"id"`
	Name      string      `json:"name"`
	Notes     string      `json:"notes"`
	StartTime time.Time   `json:"start_time"`
	EndTime   time.Time   `json:"end_time"`
	TripID    int64       `json:"trip"`
	Locations []*Location `json:"locations"`
	Version   int32       `json:"version"`
	CreatedAt time.Time   `json:"-"`
	UpdatedAt time.Time   `json:"-"`
}

type ActivityModel struct {
	DB *sql.DB
}

func (m ActivityModel) Get(id int64) (*Activity, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
    SELECT id, created_at, updated_at, name, notes, start_time, end_time trip_id, version
    FROM activities
    WHERE id = $1`

	var activity Activity

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&activity.ID,
		&activity.CreatedAt,
		&activity.UpdatedAt,
		&activity.Name,
		&activity.Notes,
		&activity.StartTime,
		&activity.EndTime,
		&activity.TripID,
		&activity.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &activity, nil
}
func (m ActivityModel) Insert(activity *Activity) error {
	query := `
    INSERT INTO activities (name, notes, start_time, end_time, trip_id)
    VALUES ($1, $2, $3, $4, $5)
    RETURNING id, created_at, version`

	args := []any{activity.Name, activity.Notes, activity.StartTime, activity.EndTime, activity.TripID}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&activity.ID, &activity.CreatedAt, &activity.Version)
}

func (m ActivityModel) GetAllByTrip(trip_id int64) ([]*Activity, error) {
	query := `
    SELECT  id, created_at, name, notes, start_time, end_time, version 
    FROM activities
    WHERE trip_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, trip_id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	activities := []*Activity{}
	// rows.Next iterates
	for rows.Next() {
		var activity Activity

		err := rows.Scan(
			&activity.ID,
			&activity.CreatedAt,
			&activity.Name,
			&activity.Notes,
			&activity.StartTime,
			&activity.EndTime,
			&activity.Version,
		)

		if err != nil {
			return nil, err
		}

		activity.TripID = trip_id
		activities = append(activities, &activity)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return activities, nil
}

func (m ActivityModel) Update(activity *Activity) error {
	query := `
    UPDATE activities
    SET name = $1, notes = $2, start_time = $3, end_time = $4
    WHERE id = $5 AND version = $6
    RETURNING version`

	args := []any{
		activity.Name,
		activity.Notes,
		activity.StartTime,
		activity.EndTime,
		activity.ID,
		activity.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&activity.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

func (m ActivityModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
    DELETE FROM activities
    WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, id)
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

func ValidateActivity(v *validator.Validator, activity *Activity) {
	// name validations
	v.Check(activity.Name != "", "name", "must be provided")
	v.Check(len(activity.Name) <= 500, "name", "must not be more than 500 bytes long")

	// notes validations
	v.Check(activity.Notes != "", "notes", "must be provided")
	v.Check(len(activity.Notes) <= 10000, "notes", "must not be more than 500 bytes long")

	// activity validations
	v.Check(activity.TripID != 0, "activity", "must be provided")

	// start_time validations
	v.Check(!activity.StartTime.IsZero(), "start_time", "must be provided")
	v.Check(!activity.StartTime.Before(time.Now()), "start_time", "must be in the future")
	v.Check(activity.StartTime.Before(activity.EndTime), "start_time", "must be before end time")

	// end_time validations
	v.Check(!activity.EndTime.IsZero(), "end_time", "must be provided")
	v.Check(!activity.EndTime.Before(time.Now()), "end_time", "must be in the future")
	v.Check(activity.EndTime.After(activity.StartTime), "end_time", "must be after start time")
}
