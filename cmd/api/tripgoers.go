package main

import (
	"errors"
	"net/http"

	"github.com/rytwalker/kagubird-api/internal/data"
)

func (app *application) addTripGoer(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email  string `json:"email"`
		TripID int64  `json:"trip"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// look up email
	user, err := app.models.Users.GetByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.invalidCredentialsResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	tripgoer := &data.TripGoer{
		UserID: user.ID,
		TripID: input.TripID,
	}

	err = app.models.TripGoers.Insert(user.ID, input.TripID)
	if err != nil {
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"tripgoer": tripgoer}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
