package main

import (
	"net/http"
	"prod/internal/handlers"
	"prod/internal/nft"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {

	go nft.UpdateCandleData(
		nft.Market_Makers_CollectionAddress,
		"Market_Makers_floor_price_data.json",
		"Market_Makers_candle_data_5min.json",
		"Market_Makers_candle_data_1hr.json",
		&nft.СandleDataMarketMakers,
		&nft.СandleDataMarketMakers,
	)
	go nft.UpdateCandleData(
		nft.Lost_Dogs_CollectionAddress,
		"Lost_Dogs_floor_price_data.json",
		"Lost_Dogs_candle_data_5min.json",
		"Lost_Dogs_candle_data_1hr.json",
		&nft.СandleDataLostDogs,
		&nft.СandleDataLostDogs,
	)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/dywetrading/getAllHistory/{address}/{timeframe}", handlers.HandleCandleData)
	r.Get("/dywetrading/getCollectionInfo/{address}", handlers.CollectionInfoHandler)

	http.ListenAndServe(":8080", r)
}
