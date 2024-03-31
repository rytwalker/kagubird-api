package main

import (
	"fmt"
	"net/http"
	"time"

	data "github.com/rytwalker/kagubird-api/internal/data"
)

func (app *application) createLocationHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Create a new location")
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
