package main

import (
	"dywego/internal/handlers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/api/v1/collection/address", handlers.GetCollectionInfo).Methods("GET")

	log.Println("Starting server on :dywe")
	log.Fatal(http.ListenAndServe(":dywe", router))
}
