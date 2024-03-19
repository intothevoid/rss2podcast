package api

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/intothevoid/rss2podcast/pkg/api/handler"
)

func StartWebService() {
	router := mux.NewRouter()

	router.HandleFunc("/generate/{topic}", handler.GenerateHandler).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", router))
}
