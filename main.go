package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/Warren-Wang-OG/go-social-media-backend/database"
)

type errorBody struct {
	Error string `json:"error"`
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK,
		database.User{
			Email: "test@example.com",
		})
}

func testErrHandler(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, http.StatusInternalServerError, errors.New("error handler default response"))
}

// wrapper for respondWithJSON that deals with errors ? (still not sure what this does tbh)
func respondWithError(w http.ResponseWriter, code int, err error) {
	respondWithJSON(w, code, errorBody{Error: err.Error()})
}

// handles http requests and return json
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	response, err := json.Marshal(payload)
	if err != nil {
		log.Fatal(err)
	}
	w.Write(response)
	w.WriteHeader(code)
}

func main() {
	// allocate http request multiplexer
	serveMux := http.NewServeMux()

	// handler to register at the "/" root path
	serveMux.HandleFunc("/", testHandler)

	// handler to register at the "/error" path
	serveMux.HandleFunc("/err", testErrHandler)

	// http server
	const addr = "localhost:8080"
	srv := http.Server{
		Handler:      serveMux,
		Addr:         addr,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}

	// wait and listen
	srv.ListenAndServe()
}
