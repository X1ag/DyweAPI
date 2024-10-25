package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"prod/internal/nft"

	"github.com/go-chi/chi/v5"
)

func CollectionInfoHandler(w http.ResponseWriter, r *http.Request) {
	address := chi.URLParam(r, "address")
	collectionData, err := nft.GetCollectionData(address)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving collection data: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(collectionData.Metadata); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode JSON response: %v", err), http.StatusInternalServerError)
	}
}
