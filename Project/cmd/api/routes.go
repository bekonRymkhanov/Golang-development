package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/Episodes", app.createEpisodeHandler)
	router.HandlerFunc(http.MethodGet, "/Episodes", app.listEpisodesHandler)
	router.HandlerFunc(http.MethodGet, "/Episodes/:id", app.showEpisodeHandler)
	router.HandlerFunc(http.MethodPatch, "/Episodes/:id", app.updateEpisodeHandler)
	router.HandlerFunc(http.MethodDelete, "/Episodes/:id", app.deleteEpisodeHandler)

	return app.recoverPanic(app.rateLimit(router))

}
