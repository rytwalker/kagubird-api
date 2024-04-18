package data

import (
	"context"
	"database/sql"
	"time"
)

type TripGoer struct {
	UserID int64 `json:"user_id"`
	TripID int64 `json:"trip_id"`
}

type TripGoerModel struct {
	DB *sql.DB
}

func (m TripGoerModel) Insert(userID int64, tripID int64) error {
	query := `
    INSERT INTO trip_goers (user_id, trip_id) 
    VALUES ($1, $2)`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, userID, tripID)
	return err
}
