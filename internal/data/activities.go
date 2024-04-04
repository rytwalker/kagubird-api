package data

import (
	"context"
	"database/sql"
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

func ValidateActivity(v *validator.Validator, activity *Activity) {
	// name validations
	v.Check(activity.Name != "", "name", "must be provided")
	v.Check(len(activity.Name) <= 500, "name", "must not be more than 500 bytes long")

	// notes validations
	v.Check(activity.Notes != "", "notes", "must be provided")
	v.Check(len(activity.Notes) <= 10000, "notes", "must not be more than 500 bytes long")

	// trip_id validations
	v.Check(activity.TripID != 0, "trip_id", "must be provided")

	// start_time validations
	v.Check(!activity.StartTime.IsZero(), "start_time", "must be provided")
	v.Check(!activity.StartTime.Before(time.Now()), "start_time", "must be in the future")
	v.Check(activity.StartTime.Before(activity.EndTime), "start_time", "must be before end time")

	// end_time validations
	v.Check(!activity.EndTime.IsZero(), "end_time", "must be provided")
	v.Check(!activity.EndTime.Before(time.Now()), "end_time", "must be in the future")
	v.Check(activity.EndTime.After(activity.StartTime), "end_time", "must be after start time")
}
