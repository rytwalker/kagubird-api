package main

import (
	"expvar"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	// ACTIVITIES
	router.HandlerFunc(http.MethodPost, "/v1/activities", app.createActivityHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/activities/:id", app.updateActivityHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/activities/:id", app.deleteActivityHandler)
	router.HandlerFunc(http.MethodGet, "/v1/activities/trip/:id", app.listActivitiesHandler)

	// LOCATIONS
	router.HandlerFunc(http.MethodPost, "/v1/locations", app.createLocationHandler)
	router.HandlerFunc(http.MethodGet, "/v1/locations/:id", app.showLocationHandler)
	// more todo...

	// METRICS
	router.Handler(http.MethodGet, "/v1/metrics", expvar.Handler())

	// STAYS
	router.HandlerFunc(http.MethodPost, "/v1/stays", app.createStayHandler)

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

	return app.metrics(app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(router)))))
}
