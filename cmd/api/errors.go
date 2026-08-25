package main

import (
	"log"
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("internal error:%s path: %s,error: %s\n", r.Method, r.URL.Path, err)
	if err := writeJSONError(w, http.StatusInternalServerError, " the server has encountered a problem."); err != nil {
		log.Printf("error inside internalServerError func: %v", err)
	}
}

func (app *application) badRequestError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("bad request error:%s path: %s,error: %s\n", r.Method, r.URL.Path, err)
	if err := writeJSONError(w, http.StatusBadRequest, err.Error()); err != nil {
		log.Printf("error inside BadRequestError func: %v", err)
	}
}

func (app *application) notFoundError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("not found error:%s path: %s,error: %s\n", r.Method, r.URL.Path, err)
	if err := writeJSONError(w, http.StatusNotFound, "resource not found"); err != nil {
		log.Printf("error inside StatusNotFound func: %v", err)
	}
}

func (app *application) conflictError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("conflict error:%s path: %s,error: %s\n", r.Method, r.URL.Path, err)
	if err := writeJSONError(w, http.StatusConflict, "resource not found"); err != nil {
		log.Printf("error inside conflictError func: %v", err)
	}
}
