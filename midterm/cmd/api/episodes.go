package main

import (
	"errors"
	"fmt"
	"net/http"
	"series.bekarysrymkhanov.net/internal/data"
	"series.bekarysrymkhanov.net/internal/validator"
)

func (app *application) createEpisodeHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title      string       `json:"title"`
		Year       int32        `json:"year"`
		Runtime    data.Runtime `json:"runtime"`
		Characters []string     `json:"characters"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	episode := &data.Episode{
		Title:      input.Title,
		Year:       input.Year,
		Runtime:    input.Runtime,
		Characters: input.Characters,
	}

	if data.ValidateMovie(v, episode); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	err = app.models.Movies.Insert(episode)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/movies/%d", episode.ID))
	err = app.writeJSON(w, http.StatusCreated, envelope{"episode": episode}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) showEpisodeHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(w, r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	episode, err := app.models.Movies.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"episode": episode}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)

	}
}
func (app *application) updateEpisodeHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(w, r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	episode, err := app.models.Movies.Get(id)
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
		Title      string       `json:"title"`
		Year       int32        `json:"year"`
		Runtime    data.Runtime `json:"runtime"`
		Characters []string     `json:"characters"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	episode.Title = input.Title
	episode.Year = input.Year
	episode.Runtime = input.Runtime
	episode.Characters = input.Characters

	v := validator.New()

	if data.ValidateMovie(v, episode); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Movies.Update(episode)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"episode": episode}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) deleteEpisodeHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(w, r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Movies.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "movie successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
