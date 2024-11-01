package nft

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
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

func UpdateCandleData(address, floorPriceFile, candleFile5Min, candleFile1Hr string, candleData5Min, candleData1Hr *CandleData) {
	var openPrice5Min, openPrice1Hr, closePrice5Min, closePrice1Hr float64
	var startTime5Min, startTime1Hr time.Time

	go func() {
		for {
			floorPrice, err := GetNFTCollectionFloor(address)
			if err != nil {
				log.Printf("Ошибка при получении floor price: %v", err)
			} else {
				// Update open and close prices for 5-minute data
				if openPrice5Min == 0 {
					openPrice5Min = floorPrice
					startTime5Min = time.Now()
				}
				closePrice5Min = floorPrice

				// Update open and close prices for 1-hour data
				if openPrice1Hr == 0 {
					openPrice1Hr = floorPrice
					startTime1Hr = time.Now()
				}
				closePrice1Hr = floorPrice

				if err := WriteFloorToFile(floorPrice, address, floorPriceFile); err != nil {
					log.Printf("Ошибка при записи в файл: %v", err)
				}
			}
			time.Sleep(20 * time.Second)
		}
	}()

	go func() {
		for {
			time.Sleep(5 * time.Minute)

			minPrice, maxPrice, err := GetCandleInfo(floorPriceFile)
			if err != nil {
				log.Printf("Ошибка при нахождении min и max значений: %v", err)
				continue
			}

			*candleData5Min = CandleData{
				StartTime: fmt.Sprintf("%d", startTime5Min.Unix()),
				EndTime:   fmt.Sprintf("%d", time.Now().Unix()),
				LowPrice:  minPrice,
				HighPrice: maxPrice,
				Open:      openPrice5Min,
				Close:     closePrice5Min,
			}

			if err := WriteCandleToFile(*candleData5Min, candleFile5Min); err != nil {
				log.Printf("Ошибка при записи данных 5-минутной свечи в файл: %v", err)
			}

			openPrice5Min = 0 // Reset for the next interval
		}
	}()

	go func() {
		for {
			time.Sleep(1 * time.Hour)

			minPrice, maxPrice, err := GetCandleInfo(floorPriceFile)
			if err != nil {
				log.Printf("Ошибка при нахождении min и max значений: %v", err)
				continue
			}

			*candleData1Hr = CandleData{
				StartTime: fmt.Sprintf("%d", startTime1Hr.Unix()),
				EndTime:   fmt.Sprintf("%d", time.Now().Unix()),
				LowPrice:  minPrice,
				HighPrice: maxPrice,
				Open:      openPrice1Hr,
				Close:     closePrice1Hr,
			}

			if err := WriteCandleToFile(*candleData1Hr, candleFile1Hr); err != nil {
				log.Printf("Ошибка при записи данных часовой свечи в файл: %v", err)
			}

			openPrice1Hr = 0 // Reset for the next interval
		}
	}()
}

func WriteCandleToFile(candle CandleData, fileName string) error {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(candle); err != nil {
		return err
	}

	return nil
}

func GetCandleInfo(fileName string) (float64, float64, error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY, 0644)
	if err != nil {
		return 0, 0, err
	}
	defer file.Close()

	var data []map[string]interface{}
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return 0, 0, err
	}

	if len(data) == 0 {
		return 0, 0, fmt.Errorf("данные отсутствуют")
	}

	minPrice := data[0]["floorPrice"].(float64)
	maxPrice := data[0]["floorPrice"].(float64)

	for _, entry := range data {
		price := entry["floorPrice"].(float64)
		if price < minPrice {
			minPrice = price
		}
		if price > maxPrice {
			maxPrice = price
		}
	}

	return minPrice, maxPrice, nil
}

func ReadLastCandleFromFile(fileName string) (CandleData, error) {
	var candle CandleData
	file, err := os.Open(fileName)
	if err != nil {
		return candle, err
	}
	defer file.Close()

	// Чтение всех данных файла и извлечение последней записи
	var candles []CandleData
	decoder := json.NewDecoder(file)
	for decoder.More() {
		var temp CandleData
		if err := decoder.Decode(&temp); err != nil {
			return candle, err
		}
		candles = append(candles, temp)
	}

	if len(candles) == 0 {
		return candle, fmt.Errorf("no candle data found in file")
	}
	candle = candles[len(candles)-1]
	return candle, nil
}
