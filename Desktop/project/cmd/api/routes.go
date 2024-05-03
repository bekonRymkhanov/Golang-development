package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/Episodes", app.requirePermission("movies:read", app.listEpisodesHandler))
	router.HandlerFunc(http.MethodPost, "/Episodes", app.requirePermission("movies:write", app.createEpisodeHandler))
	router.HandlerFunc(http.MethodGet, "/Episodes/:id", app.requirePermission("movies:read", app.showEpisodeHandler))
	router.HandlerFunc(http.MethodPatch, "/Episodes/:id", app.requirePermission("movies:write", app.updateEpisodeHandler))
	router.HandlerFunc(http.MethodDelete, "/Episodes/:id", app.requirePermission("movies:write", app.deleteEpisodeHandler))

	//router.HandlerFunc(http.MethodGet, "/healthcheck", app.healthcheckHandler)
	// router.HandlerFunc(http.MethodPost, "/Episodes", app.createEpisodeHandler)
	// router.HandlerFunc(http.MethodGet, "/Episodes", app.listEpisodesHandler)
	// router.HandlerFunc(http.MethodGet, "/Episodes/:id", app.showEpisodeHandler)
	// router.HandlerFunc(http.MethodPatch, "/Episodes/:id", app.updateEpisodeHandler)
	// router.HandlerFunc(http.MethodDelete, "/Episodes/:id", app.deleteEpisodeHandler)

	router.HandlerFunc(http.MethodPost, "/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/users/activated", app.activateUserHandler)
	router.HandlerFunc(http.MethodPost, "/tokens/authentication", app.createAuthenticationTokenHandler)

	return app.recoverPanic(app.rateLimit(app.authenticate(router)))

}
