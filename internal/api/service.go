package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/intothevoid/rss2podcast/internal/api/handler"
)

func StartWebService() {
	router := mux.NewRouter()

	// Route for generating a podcast
	router.HandleFunc("/generate/{topic}", handler.GenerateHandler).Methods("GET")

	// Route for setting the configuration of the application i.e config.yaml
	router.HandleFunc("/configure/", handler.ConfigureHandler).Methods("POST")

	// Test handler
	router.HandleFunc("/test/", TestHandler).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", router))
}

func TestHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Test handler reached")
}
