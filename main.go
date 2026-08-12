package main

import (
	"encoding/json"
	"net/http"
)

type api struct {
	addr string
}

func (a *api) getUsersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func (a *api) createUsersHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Created user ..."))
}

func main() {
	s := &api{addr: ":8080"}
	mux := http.NewServeMux()
	srv := &http.Server{
		Addr:    s.addr,
		Handler: mux,
	}
	mux.HandleFunc("GET /users", s.getUsersHandler)
	mux.HandleFunc("POST /users", s.createUsersHandler)
	srv.ListenAndServe()
}
