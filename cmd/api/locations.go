package main

import (
	"fmt"
	"net/http"
	"time"

	data "github.com/rytwalker/kagubird-api/internal/data"
	"github.com/rytwalker/kagubird-api/internal/validator"
)

func (app *application) createLocationHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name          string  `json:"name"`
		Address       string  `json:"address"`
		Lat           float64 `json:"lat"`
		Lng           float64 `json:"lng"`
		GooglePlaceID string  `json:"google_place_id"`
		Website       string  `json:"website"`
		Phone         string  `json:"phone"`
		ActivityID    int64   `json:"activity"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	location := &data.Location{
		Name:          input.Name,
		Address:       input.Address,
		Lat:           input.Lat,
		Lng:           input.Lng,
		GooglePlaceID: input.GooglePlaceID,
		Website:       input.Website,
		Phone:         input.Phone,
		ActivityID:    input.ActivityID,
	}

	v := validator.New()
	if data.ValidateLocation(v, location); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Locations.Insert(location)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/location/%d", location.ID))

	// write a json response with a 201 created status code
	err = app.writeJSON(w, http.StatusCreated, envelope{"location": location}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showLocationHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	location := data.Location{
		ID:        id,
		CreatedAt: time.Now(),
		Name:      "The BeerHive",
		Address:   "2117 Penn Ave, Pittsburgh, PA 15222",
		Lat:       40.4520999,
		Lng:       -79.9852922,
		Website:   "https://thebeerhive.com",
		Phone:     "+14129044502",
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"location": location}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
