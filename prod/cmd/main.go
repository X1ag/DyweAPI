package main

import (
	"net/http"
	"prod/internal/handlers"
	"prod/internal/nft"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {

	go nft.UpdateRealTimeCandle(
		nft.Telegram_Usernames_CollectionAddress,
		"telegramUsernamesFloorPriceArray5m",
		&nft.СandleDataTelegramUsernames,
	)

	go nft.UpdateCandleData(
		nft.Telegram_Usernames_CollectionAddress,
		"telegramUsernamesFloorPriceArray5m",
		"telegramUsernamesFloorPriceArray1h",
		&nft.СandleDataTelegramUsernames,
		&nft.СandleDataTelegramUsernames,
	)
	go nft.UpdateCandleData(
		nft.Anonymous_Telegram_Numbers_CollectionAddress,
		"anonymousTelegramNumbersPriceArray5m",
		"anonymousTelegramNumbersPriceArray1h",
		&nft.СandleDataAnonymousTelegramNumbers,
		&nft.СandleDataAnonymousTelegramNumbers,
	)
	go nft.UpdateCandleData(
		nft.TON_DNSDomains_CollectionAddress,
		"tONDNSDomainsPriceArray5m",
		"tONDNSDomainsPriceArray1h",
		&nft.СandleDataTONDNSDomains,
		&nft.СandleDataTONDNSDomains,
	)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/dywetrading/getAllHistory/{address}/{timeframe}", handlers.HandleCandleData)
	r.Get("/dywetrading/getCollectionInfo/{address}", handlers.CollectionInfoHandler)
	r.Get("/dywetrading/getFloor/{address}", handlers.CollectiongetFloor)
	r.Get("/dywetrading/ws", handlers.HandleWebSocketCandleData)

	http.ListenAndServe(":5000", r)
}
