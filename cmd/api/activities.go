package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/rytwalker/kagubird-api/internal/data"
	"github.com/rytwalker/kagubird-api/internal/validator"
)

func (app *application) createActivityHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name      string    `json:"name"`
		Notes     string    `json:"notes"`
		StartTime time.Time `json:"start_time"`
		EndTime   time.Time `json:"end_time"`
		TripID    int64     `json:"trip"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	activity := &data.Activity{
		Name:      input.Name,
		Notes:     input.Notes,
		StartTime: input.StartTime,
		EndTime:   input.EndTime,
		TripID:    input.TripID,
	}

	v := validator.New()
	if data.ValidateActivity(v, activity); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Activities.Insert(activity)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/activity/%d", activity.ID))

	// write a json response with a 201 created status code
	err = app.writeJSON(w, http.StatusCreated, envelope{"activity": activity}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
