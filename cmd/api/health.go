package main

import (
	"net/http"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]string)
	data["status"] = "ok"
	data["env"] = app.config.env
	data["version"] = version
	if err := writeJSON(w, http.StatusOK, data); err != nil {
		_ = writeJSONError(w, http.StatusInternalServerError, err.Error())
	}
}
