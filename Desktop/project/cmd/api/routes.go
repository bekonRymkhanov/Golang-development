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
	router.HandlerFunc(http.MethodGet, "/Episodes/:id/Characters", app.requirePermission("movies:write", app.showCharactersByEpisodesHandler))

	router.HandlerFunc(http.MethodGet, "/Characters", app.requirePermission("movies:read", app.listCharactersHandler))
	router.HandlerFunc(http.MethodPost, "/Characters", app.requirePermission("movies:write", app.createCharacterHandler))
	router.HandlerFunc(http.MethodGet, "/Characters/:id", app.requirePermission("movies:read", app.showCharacterHandler))
	router.HandlerFunc(http.MethodPatch, "/Characters/:id", app.requirePermission("movies:write", app.updateCharacterHandler))
	router.HandlerFunc(http.MethodDelete, "/Characters/:id", app.requirePermission("movies:write", app.deleteCharacterHandler))

	router.HandlerFunc(http.MethodGet, "/Like", app.requirePermission("movies:read", app.listLikeHandler))
	router.HandlerFunc(http.MethodPost, "/Like", app.requirePermission("movies:write", app.createLikeCommentHandler))
	router.HandlerFunc(http.MethodGet, "/Like/:id", app.requirePermission("movies:read", app.showLikeHandler))
	router.HandlerFunc(http.MethodPatch, "/Like/:id", app.requirePermission("movies:write", app.updateLikeCommentHandler))
	router.HandlerFunc(http.MethodDelete, "/Like/:id", app.requirePermission("movies:write", app.deleteLikeCommentHandler))


	router.HandlerFunc(http.MethodGet, "/healthcheck", app.healthcheckHandler)

	router.HandlerFunc(http.MethodPost, "/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/users/activated", app.activateUserHandler)
	router.HandlerFunc(http.MethodPost, "/tokens/authentication", app.createAuthenticationTokenHandler)

	return app.recoverPanic(app.rateLimit(app.authenticate(router)))

}
