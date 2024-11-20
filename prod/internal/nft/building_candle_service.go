package nft

import (
	"fmt"
	"log"
	"time"
)

var СollectionChannelstelegramUsernamesFloorPriceArray5m = make(chan CandleData, 100)
var СollectionChannelsanonymousTelegramNumbersPriceArray5m = make(chan CandleData, 100)
var СollectionChannelstONDNSDomainsPriceArray5m = make(chan CandleData, 100)

func UpdateRealTimeCandle(arrayName string, candleData *CandleData) {
	WaitUntilNextInterval(5 * time.Minute)
	time.Sleep(1 * time.Second)
	go clearChannelEvery5Minutes(СollectionChannelstelegramUsernamesFloorPriceArray5m)
	go clearChannelEvery5Minutes(СollectionChannelsanonymousTelegramNumbersPriceArray5m)
	go clearChannelEvery5Minutes(СollectionChannelstONDNSDomainsPriceArray5m)

	var floorPriceArray *[]FloorPriceData
	switch arrayName {
	case "telegramUsernamesFloorPriceArray5m":
		floorPriceArray = &telegramUsernamesFloorPriceArray5m
	case "anonymousTelegramNumbersPriceArray5m":
		floorPriceArray = &anonymousTelegramNumbersPriceArray5m
	case "tONDNSDomainsPriceArray5m":
		floorPriceArray = &tONDNSDomainsPriceArray5m
	default:
		log.Printf("Неизвестное имя массива: %s", arrayName)
		return
	}
	if len(*floorPriceArray) == 0 {
		log.Println("Нет данных для построения свечи")
		return
	}

	go func() {
		var openTime int64
		for {
			if len(*floorPriceArray) < 2 {
				openTime = time.Now().Unix()
			}

			if len(*floorPriceArray) > 0 {
				openPrice := (*floorPriceArray)[0].FloorPrice
				closePrice := (*floorPriceArray)[len(*floorPriceArray)-1].FloorPrice

				lowPrice, highPrice := getMinMax(*floorPriceArray)

				candleData.Open = openPrice
				candleData.Close = closePrice
				candleData.LowPrice = lowPrice
				candleData.HighPrice = highPrice
				candleData.OpenTime = openTime
				candleData.CloseTime = time.Now().Unix()

				if arrayName == "telegramUsernamesFloorPriceArray5m" {
					select {
					case СollectionChannelstelegramUsernamesFloorPriceArray5m <- *candleData:
					default:
						log.Printf("Канал для коллекции %s переполнен или закрыт. Размер канала: %d\n", arrayName, len(СollectionChannelstelegramUsernamesFloorPriceArray5m))
						return
					}
				} else if arrayName == "anonymousTelegramNumbersPriceArray5m" {
					select {
					case СollectionChannelsanonymousTelegramNumbersPriceArray5m <- *candleData:
					default:
						log.Printf("Канал для коллекции %s переполнен или закрыт. Размер канала: %d\n", arrayName, len(СollectionChannelsanonymousTelegramNumbersPriceArray5m))
						return
					}
				} else if arrayName == "tONDNSDomainsPriceArray5m" {
					select {
					case СollectionChannelstONDNSDomainsPriceArray5m <- *candleData:
					default:
						log.Printf("Канал для коллекции %s переполнен или закрыт. Размер канала: %d\n", arrayName, len(СollectionChannelstONDNSDomainsPriceArray5m))
						return
					}
				} else {
					log.Printf("Неизвестный тип коллекции: %s\n", arrayName)
				}
				time.Sleep(20 * time.Second)
			}
		}
	}()
}

func getMinMax(floorPrices []FloorPriceData) (minPrice, maxPrice float64) {
	minPrice, maxPrice = floorPrices[0].FloorPrice, floorPrices[0].FloorPrice

	for _, price := range floorPrices {
		if price.FloorPrice < minPrice {
			minPrice = price.FloorPrice
		}
		if price.FloorPrice > maxPrice {
			maxPrice = price.FloorPrice
		}
	}
	return
}

func (f *FloorPriceData) UnixTime() int64 {
	if f.Time == "" {
		return time.Now().Unix()
	}
	t, err := time.Parse("2006-01-02 15:04:05", f.Time)
	if err != nil {
		log.Printf("Ошибка при парсинге времени: %v", err)
		return time.Now().Unix()
	}
	return t.Unix()
}

func clearChannelEvery5Minutes(ch chan CandleData) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fmt.Printf("Очищение канала  %d\n", ch)

			for len(ch) > 0 {
				<-ch
			}

			fmt.Println("Канал очищен.")

		}
	}
}
