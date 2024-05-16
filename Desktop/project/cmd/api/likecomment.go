package main

import (
	"errors"
	"net/http"
	"fmt"
	"series.bekarysrymkhanov.net/internal/data"
	"series.bekarysrymkhanov.net/internal/validator"
)
func (app *application) createLikeCommentHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		UserID      int    `json:"user_id"`
		EpisodeID   int    `json:"episode_id"`
		CommentText string `json:"comment_text"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	likeComment := &data.LikeComment{
		UserID:      input.UserID,
		EpisodeID:   input.EpisodeID,
		CommentText: input.CommentText,
	}


	if data.ValidateLike(v, likeComment); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	err = app.models.LikeComment.Insert(likeComment)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/like/%d", likeComment.LikeID))
	err = app.writeJSON(w, http.StatusCreated, envelope{"like": likeComment}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateLikeCommentHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(w, r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	likeComment, err := app.models.LikeComment.Get(id)
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
		CommentText *string `json:"comment_text"`
		LikeCount   *int    `json:"like_count"`
		LikeID      *int    `json:"like_id"`
		Version     *int    `json:"-"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.CommentText != nil {
		likeComment.CommentText = *input.CommentText
	}
	if input.LikeCount != nil {
		likeComment.LikeCount = *input.LikeCount
	}

	v := validator.New()

	if data.ValidateLike(v, likeComment); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.LikeComment.Update(likeComment)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"likeComment": likeComment}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showLikeHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(w, r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	likeComment, err := app.models.LikeComment.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"likeComment": likeComment}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)

	}
}
func (app *application) deleteLikeCommentHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(w, r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.LikeComment.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "likeComment successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
func (app *application) listLikeHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		id int // Define LikeID as an int field
		data.Filters
	}
	v := validator.New()

	qs := r.URL.Query()

	input.id = app.readInt(qs, "id", 0, v) // Use 0 as the default value for LikeID
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)

	input.Filters.Sort = app.readString(qs, "sort", "id")

	input.Filters.SortSafelist = []string{"id", "-id"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	likeComment, metadata, err := app.models.LikeComment.GetAll(fmt.Sprintf("%d", input.id), input.Filters)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"likeComment": likeComment, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
