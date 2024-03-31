package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/rytwalker/kagubird-api/internal/data"
	"github.com/rytwalker/kagubird-api/internal/validator"
)

func (app *application) createTripHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name          string    `json:"name"`
		City          string    `json:"city"`
		StateCode     string    `json:"state_code"`
		GooglePlaceID string    `json:"google_place_id"`
		Lat           float64   `json:"lat"`
		Lng           float64   `json:"lng"`
		StartDate     time.Time `json:"start_date"`
		EndDate       time.Time `json:"end_date"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// the trip variable contains a pointer to a Trip struct.
	trip := &data.Trip{
		Name:          input.Name,
		City:          input.City,
		StateCode:     input.StateCode,
		GooglePlaceID: input.GooglePlaceID,
		Lat:           input.Lat,
		Lng:           input.Lng,
		StartDate:     input.StartDate,
		EndDate:       input.EndDate,
	}

	v := validator.New()
	if data.ValidateTrip(v, trip); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	fmt.Fprintf(w, "%+v\n", input)

	err = app.models.Trips.Insert(trip)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/trips/%d", trip.ID))

	// write a json response with a 201 created status code
	err = app.writeJSON(w, http.StatusCreated, envelope{"trip": trip}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showTripHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	trip, err := app.models.Trips.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"trip": trip}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateTripHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	trip, err := app.models.Trips.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		Name          string    `json:"name"`
		City          string    `json:"city"`
		StateCode     string    `json:"state_code"`
		GooglePlaceID string    `json:"google_place_id"`
		Lat           float64   `json:"lat"`
		Lng           float64   `json:"lng"`
		StartDate     time.Time `json:"start_date"`
		EndDate       time.Time `json:"end_date"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	trip.Name = input.Name
	trip.City = input.City
	trip.StateCode = input.StateCode
	trip.GooglePlaceID = input.GooglePlaceID
	trip.Lat = input.Lat
	trip.Lng = input.Lng
	trip.StartDate = input.StartDate
	trip.EndDate = input.EndDate

	v := validator.New()

	if data.ValidateTrip(v, trip); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Trips.Update(trip)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"trip": trip}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteTripHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Trips.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "movie successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
