package main

import (
	"encoding/json"
	"net/http"
)

// statusHandler godoc
// @Summary      Obtener estado de la API
// @Description  Retorna el estado actual de la API, incluyendo versión y entorno
// @Tags         status
// @Accept       json
// @Produce      json
// @Success      200  {object}  AppStatus
// @Router       /status [get]
func (app *application) statusHandler(w http.ResponseWriter, r *http.Request) {
	currentStatus := AppStatus{
		Status:     "Available",
		Enviroment: app.config.env,
		Version:    version,
	}
	js, err := json.MarshalIndent(currentStatus, "", "\t")
	if err != nil {
		app.logger.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}
