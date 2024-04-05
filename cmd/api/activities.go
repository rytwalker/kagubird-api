package main

import (
	"errors"
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

func (app *application) updateActivityHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	activity, err := app.models.Activities.Get(id)
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
		Name      *string    `json:"name"`
		Notes     *string    `json:"notes"`
		StartTime *time.Time `json:"start_time"`
		EndTime   *time.Time `json:"end_time"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Name != nil {
		activity.Name = *input.Name
	}

	if input.Notes != nil {
		activity.Notes = *input.Notes
	}

	if input.StartTime != nil {
		activity.StartTime = *input.StartTime
	}
	if input.EndTime != nil {
		activity.EndTime = *input.EndTime
	}

	v := validator.New()

	if data.ValidateActivity(v, activity); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Activities.Update(activity)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"activity": activity}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteActivityHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Activities.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "activity successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
