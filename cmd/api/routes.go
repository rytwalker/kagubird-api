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

	// TOKENS
	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)

	// TRIPS
	router.HandlerFunc(http.MethodPost, "/v1/trips", app.requireActivatedUser(app.createTripHandler))
	router.HandlerFunc(http.MethodGet, "/v1/trips", app.requireActivatedUser(app.listTripsHandler))
	router.HandlerFunc(http.MethodGet, "/v1/trips/:id", app.showTripHandler)
	// router.HandlerFunc(http.MethodGet, "/v1/trips/:id", app.requireActivatedUser(app.showTripHandler))
	// router.HandlerFunc(http.MethodPatch, "/v1/trips/:id", app.requireActivatedUser(app.updateTripHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/trips/:id", app.updateTripHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/trips/:id", app.requireActivatedUser(app.deleteTripHandler))

	// USERS
	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)

	return app.recoverPanic(app.rateLimit(app.authenticate(router)))
}
