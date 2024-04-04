package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Activities  ActivityModel
	Locations   LocationModel
	Permissions PermissionModel
	Tokens      TokenModel
	Trips       TripModel
	Users       UserModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Activities:  ActivityModel{DB: db},
		Locations:   LocationModel{DB: db},
		Permissions: PermissionModel{DB: db},
		Tokens:      TokenModel{DB: db},
		Trips:       TripModel{DB: db},
		Users:       UserModel{DB: db},
	}
}
