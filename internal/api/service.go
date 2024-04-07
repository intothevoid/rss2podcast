package api

import (
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

	log.Fatal(http.ListenAndServe(":8080", router))
}
