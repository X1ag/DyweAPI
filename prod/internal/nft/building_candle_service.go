package nft

import (
	"log"
	"time"
)

func UpdateRealTimeCandle(address, floorPriceArray5m string, candleData *CandleData) {
	if candleData.Open == 0 {
		candleData.Open = candleData.Close
		candleData.OpenTime = time.Now().Unix()

		candleData.LowPrice = candleData.Close
		candleData.HighPrice = candleData.Close
	}

	go func() {
		WaitUntilNextInterval(5 * time.Minute)
		for {
			floorPrice, err := GetNFTCollectionFloor(address)
			if err != nil {
				log.Printf("Ошибка при получении floor price: %v", err)
				continue
			}

			candleData.Close = floorPrice
			candleData.CloseTime = time.Now().Unix()

			if floorPrice < candleData.LowPrice || candleData.LowPrice == 0 {
				candleData.LowPrice = floorPrice
			}

			if floorPrice > candleData.HighPrice {
				candleData.HighPrice = floorPrice
			}

			if err := WriteFloorToArray(floorPrice, address, floorPriceArray5m); err != nil {
				log.Printf("Ошибка при записи в массив 5-минутных данных: %v", err)
			}

			log.Printf("Текущая строящаяся свеча (каждые 20 секунд): %+v", *candleData)

			time.Sleep(20 * time.Second)
		}
	}()
}
