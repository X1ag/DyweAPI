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
		nft.Telegram_Usernames_CollectionAddress,
		"telegramUsernamesFloorPriceArray",
		"Telegram_Usernames_candle_data_5min.json",
		"Telegram_Usernames_candle_data_1hr.json",
		&nft.СandleDataTelegramUsernames,
		&nft.СandleDataTelegramUsernames,
	)
	go nft.UpdateCandleData(
		nft.Anonymous_Telegram_Numbers_CollectionAddress,
		"anonymousTelegramNumbersPriceArray",
		"Anonymous_Telegram_Numbers_data_5min.json",
		"Anonymous_Telegram_Numbers_candle_data_1hr.json",
		&nft.СandleDataAnonymousTelegramNumbers,
		&nft.СandleDataAnonymousTelegramNumbers,
	)
	go nft.UpdateCandleData(
		nft.TON_DNSDomains_CollectionAddress,
		"tONDNSDomainsPriceArray",
		"TON_DNSDomains_data_5min.json",
		"TON_DNSDomains_candle_data_1hr.json",
		&nft.СandleDataTONDNSDomains,
		&nft.СandleDataTONDNSDomains,
	)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/dywetrading/getAllHistory/{address}/{timeframe}", handlers.HandleCandleData)
	r.Get("/dywetrading/getCollectionInfo/{address}", handlers.CollectionInfoHandler)
	r.Get("/dywetrading/getFloor/{address}", handlers.CollectiongetFloor)

	http.ListenAndServe(":5000", r)
}
