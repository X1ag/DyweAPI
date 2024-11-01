package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"prod/internal/nft"

	"github.com/go-chi/chi/v5"
)

func HandleCandleData(w http.ResponseWriter, r *http.Request) {
	address := chi.URLParam(r, "address")
	timeframe := chi.URLParam(r, "timeframe")

	var candleDataResponse nft.CandleData
	var candleFile string

	if address == nft.Market_Makers_CollectionAddress {
		if timeframe == "5min" {
			candleFile = "Market_Makers_candle_data_5min.json"
		} else if timeframe == "1hr" {
			candleFile = "Market_Makers_candle_data_1hr.json"
		} else {
			http.Error(w, "Invalid timeframe", http.StatusBadRequest)
			return
		}
	} else if address == nft.Lost_Dogs_CollectionAddress {
		if timeframe == "5min" {
			candleFile = "Lost_Dogs_candle_data_5min.json"
		} else if timeframe == "1hr" {
			candleFile = "Lost_Dogs_candle_data_1hr.json"
		} else {
			http.Error(w, "Invalid timeframe", http.StatusBadRequest)
			return
		}
	} else {
		http.Error(w, "Collection not found", http.StatusNotFound)
		return
	}

	candleDataResponse, err := nft.ReadLastCandleFromFile(candleFile)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading candle data: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(candleDataResponse)
}
