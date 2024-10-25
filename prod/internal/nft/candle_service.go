package nft

import (
	"fmt"
	"log"
	"time"
)

type CandleData struct {
	StartTime string  `json:"startTime"`
	EndTime   string  `json:"endTime"`
	LowPrice  float64 `json:"lowPrice"`
	HighPrice float64 `json:"highPrice"`
	Open      float64 `json:"open"`
	Close     float64 `json:"close"`
}

var СandleDataMarketMakers CandleData
var СandleDataLostDogs CandleData

var Market_Makers_CollectionAddress = "EQBDMXqg2YcGmMnn5_bXG63y-hh_YNV0dx-ylx-vL3v_WZt4"
var Lost_Dogs_CollectionAddress = "EQAl_hUCAeEv-fKtGxYtITAS6PPxuMRaQwHj0QAHeWe6ZSD0"

func UpdateCandleData(address string, fileName string, candleData *CandleData) {
	var openPrice, closePrice float64
	var startTime time.Time

	go func() {
		for {
			floorPrice, err := GetNFTCollectionFloor(address)
			if err != nil {
				log.Printf("Ошибка при получении floor price: %v", err)
			} else {
				if openPrice == 0 {
					openPrice = floorPrice
					startTime = time.Now()
				}
				closePrice = floorPrice

				if err := WriteFloorToFile(floorPrice, address, fileName); err != nil {
					log.Printf("Ошибка при записи в файл: %v", err)
				}
			}
			time.Sleep(20 * time.Second)
		}
	}()

	go func() {
		for {
			time.Sleep(1 * time.Hour)

			minPrice, maxPrice, err := GetCandleInfo(fileName)
			if err != nil {
				log.Printf("Ошибка при нахождении min и max значений: %v", err)
				continue
			}

			*candleData = CandleData{
				StartTime: fmt.Sprintf("%d", startTime.Unix()),
				EndTime:   fmt.Sprintf("%d", time.Now().Unix()),
				LowPrice:  minPrice,
				HighPrice: maxPrice,
				Open:      openPrice,
				Close:     closePrice,
			}

			fmt.Printf("Candle data обновлены для %s: %v\n", address, *candleData)

			openPrice = 0
		}
	}()
}

func GetCandleData(address string) (CandleData, error) {
	var candleData CandleData

	if address == Market_Makers_CollectionAddress {
		candleData = СandleDataMarketMakers
	} else if address == Lost_Dogs_CollectionAddress {
		candleData = СandleDataLostDogs
	} else {
		return candleData, fmt.Errorf("collection not found")
	}

	return candleData, nil
}
