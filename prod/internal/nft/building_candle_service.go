package nft

import (
	"fmt"
	"log"
	"time"
)

var СollectionChannelstelegramUsernamesFloorPriceArray5m = make(chan CandleData, 100)
var СollectionChannelsanonymousTelegramNumbersPriceArray5m = make(chan CandleData, 100)
var СollectionChannelstONDNSDomainsPriceArray5m = make(chan CandleData, 100)
var СollectionChannelstelegramUsernamesFloorPriceArray1h = make(chan CandleData, 100)
var СollectionChannelsanonymousTelegramNumbersPriceArray1h = make(chan CandleData, 100)
var СollectionChannelstONDNSDomainsPriceArray1h = make(chan CandleData, 100)
var СollectionChannelstelegramUsernamesFloorPriceArray15m = make(chan CandleData, 100)
var СollectionChannelsanonymousTelegramNumbersPriceArray15m = make(chan CandleData, 100)
var СollectionChannelstONDNSDomainsPriceArray15m = make(chan CandleData, 100)
var СollectionChannelstelegramUsernamesFloorPriceArray30m = make(chan CandleData, 100)
var СollectionChannelsanonymousTelegramNumbersPriceArray30m = make(chan CandleData, 100)
var СollectionChannelstONDNSDomainsPriceArray30m = make(chan CandleData, 100)
var СollectionChannelstelegramUsernamesFloorPriceArray4h = make(chan CandleData, 100)
var СollectionChannelsanonymousTelegramNumbersPriceArray4h = make(chan CandleData, 100)
var СollectionChannelstONDNSDomainsPriceArray4h = make(chan CandleData, 100)

func UpdateRealTimeCandle5m(arrayName string, candleData *CandleData) {
	WaitUntilNextInterval(5 * time.Minute)
	time.Sleep(1 * time.Second)
	go clearChannel(СollectionChannelstelegramUsernamesFloorPriceArray5m, 5*time.Minute)
	go clearChannel(СollectionChannelsanonymousTelegramNumbersPriceArray5m, 5*time.Minute)
	go clearChannel(СollectionChannelstONDNSDomainsPriceArray5m, 5*time.Minute)

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

func UpdateRealTimeCandle1h(arrayName string, candleData *CandleData) {
	WaitUntilNextInterval(1 * time.Hour)
	time.Sleep(1 * time.Second)
	go clearChannel(СollectionChannelstelegramUsernamesFloorPriceArray1h, 1*time.Hour)
	go clearChannel(СollectionChannelsanonymousTelegramNumbersPriceArray1h, 1*time.Hour)
	go clearChannel(СollectionChannelstONDNSDomainsPriceArray1h, 1*time.Hour)

	var floorPriceArray *[]FloorPriceData
	switch arrayName {
	case "telegramUsernamesFloorPriceArray1h":
		floorPriceArray = &telegramUsernamesFloorPriceArray1h
	case "anonymousTelegramNumbersPriceArray1h":
		floorPriceArray = &anonymousTelegramNumbersPriceArray1h
	case "tONDNSDomainsPriceArray1h":
		floorPriceArray = &tONDNSDomainsPriceArray1h
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

				if arrayName == "telegramUsernamesFloorPriceArray1h" {
					select {
					case СollectionChannelstelegramUsernamesFloorPriceArray1h <- *candleData:
					default:
						log.Printf("Канал для коллекции %s переполнен или закрыт. Размер канала: %d\n", arrayName, len(СollectionChannelstelegramUsernamesFloorPriceArray1h))
						return
					}
				} else if arrayName == "anonymousTelegramNumbersPriceArray1h" {
					select {
					case СollectionChannelsanonymousTelegramNumbersPriceArray1h <- *candleData:
					default:
						log.Printf("Канал для коллекции %s переполнен или закрыт. Размер канала: %d\n", arrayName, len(СollectionChannelsanonymousTelegramNumbersPriceArray1h))
						return
					}
				} else if arrayName == "tONDNSDomainsPriceArray1h" {
					select {
					case СollectionChannelstONDNSDomainsPriceArray1h <- *candleData:
					default:
						log.Printf("Канал для коллекции %s переполнен или закрыт. Размер канала: %d\n", arrayName, len(СollectionChannelstONDNSDomainsPriceArray1h))
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

func UpdateRealTimeCandle15m(arrayName string, candleData *CandleData) {
	WaitUntilNextInterval(15 * time.Minute)
	time.Sleep(1 * time.Second)
	go clearChannel(СollectionChannelstelegramUsernamesFloorPriceArray15m, 15*time.Minute)
	go clearChannel(СollectionChannelsanonymousTelegramNumbersPriceArray15m, 15*time.Minute)
	go clearChannel(СollectionChannelstONDNSDomainsPriceArray15m, 15*time.Minute)

	var floorPriceArray *[]FloorPriceData
	switch arrayName {
	case "telegramUsernamesFloorPriceArray15m":
		floorPriceArray = &telegramUsernamesFloorPriceArray15m
	case "anonymousTelegramNumbersPriceArray15m":
		floorPriceArray = &anonymousTelegramNumbersPriceArray15m
	case "tONDNSDomainsPriceArray15m":
		floorPriceArray = &tONDNSDomainsPriceArray15m
	default:
		log.Printf("Неизвестное имя массива: %s", arrayName)
		return
	}
	if len(*floorPriceArray) == 0 {
		log.Println("Нет данных для построения свечи")
		return
	}

	go func() {
		for {
			time.Sleep(15 * time.Minute)
			ClearArray(arrayName)
			log.Printf("Массив %s был очищен", arrayName)
		}
	}()

	go func() {
		ClearArray(arrayName)
		var openTime int64
		for {

			if len(*floorPriceArray) < 1 {
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

				if arrayName == "telegramUsernamesFloorPriceArray15m" {
					select {
					case СollectionChannelstelegramUsernamesFloorPriceArray15m <- *candleData:
					default:
						log.Printf("Канал для коллекции %s переполнен или закрыт. Размер канала: %d\n", arrayName, len(СollectionChannelstelegramUsernamesFloorPriceArray15m))
						return
					}
				} else if arrayName == "anonymousTelegramNumbersPriceArray15m" {
					select {
					case СollectionChannelsanonymousTelegramNumbersPriceArray15m <- *candleData:
					default:
						log.Printf("Канал для коллекции %s переполнен или закрыт. Размер канала: %d\n", arrayName, len(СollectionChannelsanonymousTelegramNumbersPriceArray15m))
						return
					}
				} else if arrayName == "tONDNSDomainsPriceArray15m" {
					select {
					case СollectionChannelstONDNSDomainsPriceArray15m <- *candleData:
					default:
						log.Printf("Канал для коллекции %s переполнен или закрыт. Размер канала: %d\n", arrayName, len(СollectionChannelstONDNSDomainsPriceArray15m))
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

func UpdateRealTimeCandle30m(arrayName string, candleData *CandleData) {
	WaitUntilNextInterval(30 * time.Minute)
	time.Sleep(1 * time.Second)
	go clearChannel(СollectionChannelstelegramUsernamesFloorPriceArray30m, 30*time.Minute)
	go clearChannel(СollectionChannelsanonymousTelegramNumbersPriceArray30m, 30*time.Minute)
	go clearChannel(СollectionChannelstONDNSDomainsPriceArray30m, 30*time.Minute)

	var floorPriceArray *[]FloorPriceData
	switch arrayName {
	case "telegramUsernamesFloorPriceArray30m":
		floorPriceArray = &telegramUsernamesFloorPriceArray30m
	case "anonymousTelegramNumbersPriceArray30m":
		floorPriceArray = &anonymousTelegramNumbersPriceArray30m
	case "tONDNSDomainsPriceArray30m":
		floorPriceArray = &tONDNSDomainsPriceArray30m
	default:
		log.Printf("Неизвестное имя массива: %s", arrayName)
		return
	}
	if len(*floorPriceArray) == 0 {
		log.Println("Нет данных для построения свечи")
		return
	}

	go func() {
		for {
			time.Sleep(30 * time.Minute)
			ClearArray(arrayName)
			log.Printf("Массив %s был очищен", arrayName)
		}
	}()

	go func() {
		ClearArray(arrayName)
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

				if arrayName == "telegramUsernamesFloorPriceArray30m" {
					select {
					case СollectionChannelstelegramUsernamesFloorPriceArray30m <- *candleData:
					default:
						log.Printf("Канал для коллекции %s переполнен или закрыт. Размер канала: %d\n", arrayName, len(СollectionChannelstelegramUsernamesFloorPriceArray30m))
						return
					}
				} else if arrayName == "anonymousTelegramNumbersPriceArray30m" {
					select {
					case СollectionChannelsanonymousTelegramNumbersPriceArray30m <- *candleData:
					default:
						log.Printf("Канал для коллекции %s переполнен или закрыт. Размер канала: %d\n", arrayName, len(СollectionChannelsanonymousTelegramNumbersPriceArray30m))
						return
					}
				} else if arrayName == "tONDNSDomainsPriceArray30m" {
					select {
					case СollectionChannelstONDNSDomainsPriceArray30m <- *candleData:
					default:
						log.Printf("Канал для коллекции %s переполнен или закрыт. Размер канала: %d\n", arrayName, len(СollectionChannelstONDNSDomainsPriceArray30m))
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

func UpdateRealTimeCandle4h(arrayName string, candleData *CandleData) {
	WaitUntilNextInterval(4 * time.Hour)
	time.Sleep(1 * time.Second)
	go clearChannel(СollectionChannelstelegramUsernamesFloorPriceArray4h, 4*time.Hour)
	go clearChannel(СollectionChannelsanonymousTelegramNumbersPriceArray4h, 4*time.Hour)
	go clearChannel(СollectionChannelstONDNSDomainsPriceArray4h, 4*time.Hour)

	var floorPriceArray *[]FloorPriceData
	switch arrayName {
	case "telegramUsernamesFloorPriceArray4h":
		floorPriceArray = &telegramUsernamesFloorPriceArray4h
	case "anonymousTelegramNumbersPriceArray4h":
		floorPriceArray = &anonymousTelegramNumbersPriceArray4h
	case "tONDNSDomainsPriceArray4h":
		floorPriceArray = &tONDNSDomainsPriceArray4h
	default:
		log.Printf("Неизвестное имя массива: %s", arrayName)
		return
	}
	if len(*floorPriceArray) == 0 {
		log.Println("Нет данных для построения свечи")
		return
	}

	go func() {
		for {
			time.Sleep(4 * time.Hour)
			ClearArray(arrayName)
			log.Printf("Массив %s был очищен", arrayName)
		}
	}()

	go func() {
		ClearArray(arrayName)
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

				if arrayName == "telegramUsernamesFloorPriceArray4h" {
					select {
					case СollectionChannelstelegramUsernamesFloorPriceArray4h <- *candleData:
					default:
						log.Printf("Канал для коллекции %s переполнен или закрыт. Размер канала: %d\n", arrayName, len(СollectionChannelstelegramUsernamesFloorPriceArray4h))
						return
					}
				} else if arrayName == "anonymousTelegramNumbersPriceArray4h" {
					select {
					case СollectionChannelsanonymousTelegramNumbersPriceArray4h <- *candleData:
					default:
						log.Printf("Канал для коллекции %s переполнен или закрыт. Размер канала: %d\n", arrayName, len(СollectionChannelsanonymousTelegramNumbersPriceArray4h))
						return
					}
				} else if arrayName == "tONDNSDomainsPriceArray4h" {
					select {
					case СollectionChannelstONDNSDomainsPriceArray4h <- *candleData:
					default:
						log.Printf("Канал для коллекции %s переполнен или закрыт. Размер канала: %d\n", arrayName, len(СollectionChannelstONDNSDomainsPriceArray4h))
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

func clearChannel(ch chan CandleData, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fmt.Printf("Очищение канала через %v\n", interval)
			for len(ch) > 0 {
				<-ch
			}

			fmt.Println("Канал очищен.")
		}
	}
}
