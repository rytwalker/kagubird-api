package main

import (
	"fmt"
	"net/http"
	"time"

	data "github.com/rytwalker/kagubird-api/internal/data"
	"github.com/rytwalker/kagubird-api/internal/validator"
)

func (app *application) createStayHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name      string    `json:"name"`
		Address   string    `json:"address"`
		Lat       float64   `json:"lat"`
		Lng       float64   `json:"lng"`
		StartTime time.Time `json:"start_time"`
		EndTime   time.Time `json:"end_time"`
		Link      string    `json:"link"`
		Phone     string    `json:"phone"`
		Type      string    `json:"type"`
		TripID    int64     `json:"trip"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	stay := &data.Stay{
		Name:      input.Name,
		Address:   input.Address,
		Lat:       input.Lat,
		Lng:       input.Lng,
		StartTime: input.StartTime,
		EndTime:   input.EndTime,
		Link:      input.Link,
		Phone:     input.Phone,
		Type:      input.Type,
		TripID:    input.TripID,
	}

	v := validator.New()
	if data.ValidateStay(v, stay); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Stays.Insert(stay)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/stay/%d", stay.ID))

	// write a json response with a 201 created status code
	err = app.writeJSON(w, http.StatusCreated, envelope{"stay": stay}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
