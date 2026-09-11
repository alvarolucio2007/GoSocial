package main

import (
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("internal error", "method", r.Method, "path", r.URL.Path, "error", err)
	if err := writeJSONError(w, http.StatusInternalServerError, " the server has encountered a problem."); err != nil {
		app.logger.Errorf("error while attempting function writeJSONError inside internalServerError %v", err)
	}
}

func (app *application) badRequestError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnw("bad request error", "method", r.Method, "path", r.URL.Path, "error", err)
	if err := writeJSONError(w, http.StatusBadRequest, err.Error()); err != nil {
		app.logger.Errorf("error while attempting function writeJSONError inside badRequestError %v", err)
	}
}

func (app *application) notFoundError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnw("not found error", "method", r.Method, "path", r.URL.Path, "error", err)
	if err := writeJSONError(w, http.StatusNotFound, "resource not found"); err != nil {
		app.logger.Errorf("error while attempting function writeJSONError inside notFoundError %v", err)
	}
}

func (app *application) conflictError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnw("conflict error", "method", r.Method, "path", r.URL.Path, "error", err)
	if err := writeJSONError(w, http.StatusConflict, "resource not found"); err != nil {
		app.logger.Errorf("error while attempting function writeJSONError inside conflictError %v", err)
	}
}

func (app *application) unauthorizedBasicError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnw("unauthorized basic error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
	if err := writeJSONError(w, http.StatusUnauthorized, "unauthorized basic"); err != nil {
		app.logger.Errorf("error while attempting function writeJSONError inside unauthorizedError %v", err)
	}
}

func (app *application) unauthorizedError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnw("unauthorized error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	if err := writeJSONError(w, http.StatusUnauthorized, "unauthorized"); err != nil {
		app.logger.Errorf("error while attempting function writeJSONError inside unauthorizedError %v", err)
	}
}
