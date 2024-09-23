package handlers

import (
	"encoding/json"
	"net/http"

	"dywego/internal/services" // должно коннектить

	"github.com/gorilla/mux"
)

// GetCollectionInfo обрабатывает запрос для получения информации о коллекции по address
func GetCollectionInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	collectionID := vars["address"]

	collection, err := services.GetCollection(collectionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(collection)
}
