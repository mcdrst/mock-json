package main

import (
	"log"
	"net/http"
)

// application errors
func serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Println(err)
	message := "the server encountered a problem an could not process your request"
	errorResponse(w, r, http.StatusInternalServerError, message)
}

func errorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	env := mapStringInterface{"error": message}
	err := writeJSON(w, status, env)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
	}
}

func notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the request resource could not be found"
	errorResponse(w, r, http.StatusNotFound, message)
}

func failedValidationResponse(w http.ResponseWriter, r *http.Request, errors mapStringInterface) {
	errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

func parseFloatError(w http.ResponseWriter, r *http.Request) {
	message := "numeric parameter conversion error"
	errorResponse(w, r, http.StatusInternalServerError, message)
}

func parseJSONError(w http.ResponseWriter, r *http.Request) {
	message := "request body data error"
	errorResponse(w, r, http.StatusInternalServerError, message)
}
