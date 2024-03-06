package main

import (
	"encoding/json"
	"net/http"
)

func (app *application) writeJSON(w http.ResponseWriter, status int, data interface{}, wrap string) error {
	wrapper := make(map[string]interface{})

	wrapper[wrap] = data
	js, err := json.Marshal(wrapper)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	//w.WriteHeader(status)
	w.Write(js)
	return nil
}

func (app *application) errJSON(w http.ResponseWriter, err error) {
	type jsonerror struct {
		Message string `json:"message"`
	}
	theError := jsonerror{Message: err.Error()}
	app.writeJSON(w, http.StatusBadRequest, theError, "error")
}
