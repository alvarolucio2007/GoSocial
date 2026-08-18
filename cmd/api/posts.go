package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/alvarolucio2007/GoSocial/internal/store"
	"github.com/go-chi/chi/v5"
)

type CreatePostPayload struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
	UserID  int      `json:"user_id"`
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreatePostPayload
	if err := readJSON(w, r, &payload); err != nil {
		_ = writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	post := &store.Post{
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
		UserID:  int64(payload.UserID),
	}
	if post.Tags == nil {
		post.Tags = []string{}
	}
	ctx := r.Context()
	if err := app.storage.Posts.Create(ctx, post); err != nil {
		_ = writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err := writeJSON(w, http.StatusCreated, post); err != nil {
		_ = writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (app *application) readPostHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "postID")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		_ = writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	ctx := r.Context()
	post, err := app.storage.Posts.Read(ctx, int(id))
	if err != nil {
		switch {
		case errors.Is(err, store.ErrPostNotFound):
			_ = writeJSONError(w, http.StatusNotFound, err.Error())
		default:
			_ = writeJSONError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	if err := writeJSON(w, http.StatusCreated, &post); err != nil {
		_ = writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (app *application) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "postID")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		_ = writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	var payload store.Post
	if err := readJSON(w, r, &payload); err != nil {
		_ = writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	post := &store.Post{
		ID:      id,
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
	}
	ctx := r.Context()
	if err := app.storage.Posts.Update(ctx, post); err != nil {
		switch {
		case errors.Is(err, store.ErrPostNotFound):
			_ = writeJSONError(w, http.StatusNotFound, err.Error())
		default:
			_ = writeJSONError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	if err := writeJSON(w, http.StatusOK, &post); err != nil {
		_ = writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "postID")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		_ = writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	ctx := r.Context()
	if err := app.storage.Posts.Delete(ctx, int(id)); err != nil {
		switch {
		case errors.Is(err, store.ErrPostNotFound):
			_ = writeJSONError(w, http.StatusNotFound, err.Error())
		default:
			_ = writeJSONError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
