package main

import (
	"net/http"
	"time"
)

func testHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write([]byte("{}"))
}

func main() {
	// allocate http request multiplexer
	serveMux := http.NewServeMux()

	// handler to register at the "/" root path
	serveMux.HandleFunc("/", testHandler)

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
