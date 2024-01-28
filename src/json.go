package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func parseReqBody(r *http.Request, dest interface{}) error {
	decoder := json.NewDecoder(r.Body)
	return decoder.Decode(dest)
}

func respondWithJson(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Error marshalling JSON: %s\n", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(statusCode)
	w.Write(data)
}

func respondWithError(w http.ResponseWriter, statusCode int, msg string) {
	type errorResponse struct {
		Error string `json:"error"`
	}
	respondWithJson(w, statusCode, errorResponse{Error: msg})
}
