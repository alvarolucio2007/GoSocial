package main

import (
	"log"
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("internal error:%s path: %s,error: %s\n", r.Method, r.URL.Path, err)
	writeJSONError(w, http.StatusInternalServerError, " the server has encountered a problem.")
}

func (app *application) badRequestError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("bad request error:%s path: %s,error: %s\n", r.Method, r.URL.Path, err)
	writeJSONError(w, http.StatusInternalServerError, err.Error())
}

func (app *application) notFoundError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("not found error:%s path: %s,error: %s\n", r.Method, r.URL.Path, err)
	writeJSONError(w, http.StatusNotFound, "resource not found")
}
