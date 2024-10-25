package handlers

import (
	"encoding/json"
	"net/http"
	"prod/internal/nft"

	"github.com/go-chi/chi/v5"
)

func HandleCandleData(w http.ResponseWriter, r *http.Request) {
	address := chi.URLParam(r, "address")

	var candleDataResponse nft.CandleData

	if address == nft.Market_Makers_CollectionAddress {
		candleDataResponse = nft.СandleDataMarketMakers
	} else if address == nft.Lost_Dogs_CollectionAddress {
		candleDataResponse = nft.СandleDataLostDogs
	} else {
		http.Error(w, "Collection not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(candleDataResponse)
}
