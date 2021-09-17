package web

import (
	"encoding/json"
	"net/http"
)

func Respond(w http.ResponseWriter, data interface{}, statusCode int) error {
	w.WriteHeader(statusCode)
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	if statusCode != 200 {
		w.WriteHeader(statusCode)
	}
	if _, err := w.Write(payload); err != nil {
		return err
	}

	return nil
}

