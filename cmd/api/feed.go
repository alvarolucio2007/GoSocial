package main

import (
	"net/http"

	"github.com/alvarolucio2007/GoSocial/internal/store"
)

func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	fq := store.PaginatedFeedQuery{
		Limit:  20,
		Offset: 0,
		Sort:   "desc",
	}
	fq, err := fq.Parse(r)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}
	if err := Validate.Struct(fq); err != nil {
		app.badRequestError(w, r, err)
		return
	}
	ctx := r.Context()
	feed, err := app.storage.Posts.GetUserFeed(ctx, int64(10), fq)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	if err := app.jsonResponse(w, http.StatusOK, feed); err != nil {
		app.internalServerError(w, r, err)
	}
}
