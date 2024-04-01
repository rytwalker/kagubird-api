package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	// ITINERY ITEMS
	// todo...

	// LOCATIONS
	router.HandlerFunc(http.MethodPost, "/v1/locations", app.createLocationHandler)
	router.HandlerFunc(http.MethodGet, "/v1/locations/:id", app.showLocationHandler)
	// more todo...

	// TRIPS
	router.HandlerFunc(http.MethodPost, "/v1/trips", app.createTripHandler)
	router.HandlerFunc(http.MethodGet, "/v1/trips", app.listTripsHandler)
	router.HandlerFunc(http.MethodGet, "/v1/trips/:id", app.showTripHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/trips/:id", app.updateTripHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/trips/:id", app.deleteTripHandler)

	// USERS
	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)

	return app.recoverPanic(app.rateLimit(router))
}
