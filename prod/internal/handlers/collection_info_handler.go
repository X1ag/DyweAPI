package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"prod/internal/nft"

	"github.com/go-chi/chi/v5"
)

type FloorPriceResponse struct {
	FloorPrice float64 `json:"floor_price"`
}

func CollectiongetFloor(w http.ResponseWriter, r *http.Request) {
	nftCollectionAddress := chi.URLParam(r, "address")
	if nftCollectionAddress == "" {
		http.Error(w, "Missing collection address", http.StatusBadRequest)
		return
	}
	floorPrice, err := nft.GetNFTCollectionFloor(nftCollectionAddress)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get floor price: %v", err), http.StatusInternalServerError)
		return
	}
	response := FloorPriceResponse{
		FloorPrice: floorPrice,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode JSON response: %v", err), http.StatusInternalServerError)
	}
}

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
