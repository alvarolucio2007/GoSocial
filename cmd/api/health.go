package main

import (
	"net/http"
)

// healthcheckHandler godoc
//
//	@Summary		Healthcheck
//	@Description	Healthcheck endpoint
//	@Tags			ops
//	@Produce		json
//	@Success		200	{object}	string	"ok"
//	@Router			/health [get]
func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]string)
	data["status"] = "ok"
	data["env"] = app.config.env
	data["version"] = version
	if err := writeJSON(w, http.StatusOK, data); err != nil {
		app.internalServerError(w, r, err)
	}
}
