package main

import (
	"fmt"
	"net/http"
	"prod/internal/handlers"
	"prod/internal/nft"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	go nft.UpdateRealTimeCandle5m(
		"telegramUsernamesFloorPriceArray5m",
		&nft.СandleDataTelegramUsernames,
	)

	go nft.UpdateRealTimeCandle5m(
		"anonymousTelegramNumbersPriceArray5m",
		&nft.СandleDataTelegramUsernames,
	)

	go nft.UpdateRealTimeCandle5m(
		"tONDNSDomainsPriceArray5m",
		&nft.СandleDataTelegramUsernames,
	)

	go nft.UpdateRealTimeCandle1h(
		"telegramUsernamesFloorPriceArray1h",
		&nft.СandleDataTelegramUsernames,
	)

	go nft.UpdateRealTimeCandle1h(
		"anonymousTelegramNumbersPriceArray1h",
		&nft.СandleDataTelegramUsernames,
	)

	go nft.UpdateRealTimeCandle1h(
		"tONDNSDomainsPriceArray1h",
		&nft.СandleDataTelegramUsernames,
	)

	go nft.UpdateRealTimeCandle15m(
		"telegramUsernamesFloorPriceArray15m",
		&nft.СandleDataTelegramUsernames,
	)

	go nft.UpdateRealTimeCandle15m(
		"anonymousTelegramNumbersPriceArray15m",
		&nft.СandleDataTelegramUsernames,
	)

	go nft.UpdateRealTimeCandle15m(
		"tONDNSDomainsPriceArray15m",
		&nft.СandleDataTelegramUsernames,
	)

	go nft.UpdateRealTimeCandle30m(
		"telegramUsernamesFloorPriceArray30m",
		&nft.СandleDataTelegramUsernames,
	)

	go nft.UpdateRealTimeCandle30m(
		"anonymousTelegramNumbersPriceArray30m",
		&nft.СandleDataTelegramUsernames,
	)

	go nft.UpdateRealTimeCandle30m(
		"tONDNSDomainsPriceArray30m",
		&nft.СandleDataTelegramUsernames,
	)

	go nft.UpdateCandleData(
		nft.Telegram_Usernames_CollectionAddress,
		"telegramUsernamesFloorPriceArray5m",
		"telegramUsernamesFloorPriceArray1h",
		"telegramUsernamesFloorPriceArray15m",
		"telegramUsernamesFloorPriceArray30m",
		"telegramUsernamesFloorPriceArray4h",
		&nft.СandleDataTelegramUsernames,
		&nft.СandleDataTelegramUsernames,
	)
	go nft.UpdateCandleData(
		nft.Anonymous_Telegram_Numbers_CollectionAddress,
		"anonymousTelegramNumbersPriceArray5m",
		"anonymousTelegramNumbersPriceArray1h",
		"anonymousTelegramNumbersPriceArray15m",
		"anonymousTelegramNumbersPriceArray30m",
		"anonymousTelegramNumbersPriceArray1h",
		&nft.СandleDataAnonymousTelegramNumbers,
		&nft.СandleDataAnonymousTelegramNumbers,
	)
	go nft.UpdateCandleData(
		nft.TON_DNSDomains_CollectionAddress,
		"tONDNSDomainsPriceArray5m",
		"tONDNSDomainsPriceArray1h",
		"tONDNSDomainsPriceArray15m",
		"tONDNSDomainsPriceArray30m",
		"tONDNSDomainsPriceArray1h",
		&nft.СandleDataTONDNSDomains,
		&nft.СandleDataTONDNSDomains,
	)

	fmt.Println("ОЖИДАНИЕ НАЧАЛА ПЯТИМИНУТНОГО ИНТЕРВАЛА")

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/dywetrading/getAllHistory/{address}/{timeframe}", handlers.HandleCandleData)
	r.Get("/dywetrading/getCollectionInfo/{address}", handlers.CollectionInfoHandler)
	r.Get("/dywetrading/getFloor/{address}", handlers.CollectiongetFloor)
	r.Get("/dywetrading/ws/{name}/{timeframe}", handlers.HandleWebSocketCandleData)

	http.ListenAndServe(":5000", r)

	fmt.Println("Ожидание начала пятиминутного интервала")
}
