package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/Episodes", app.createEpisodeHandler)
	router.HandlerFunc(http.MethodGet, "/Episodes/:id", app.showEpisodeHandler)
	router.HandlerFunc(http.MethodPut, "/Episodes/:id", app.updateEpisodeHandler)
	router.HandlerFunc(http.MethodDelete, "/Episodes/:id", app.deleteEpisodeHandler)

	return router

}
