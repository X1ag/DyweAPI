package nft

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"prod/db"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
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

	ctx := context.Background()

	dbPool, err := db.NewClient(ctx, 3, db.StorageConfig{})
	if err != nil {
		log.Fatalf("Ошибка при создании пула соединений: %v", err)
	}
	defer dbPool.Close()

	go func() {
		for {
			floorPrice, err := GetNFTCollectionFloor(address)
			if err != nil {
				log.Printf("Ошибка при получении floor price: %v", err)
			} else {
				if openPrice5Min == 0 {
					openPrice5Min = floorPrice
					startTime5Min = time.Now()
				}
				closePrice5Min = floorPrice

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

			if err := WriteCandleToDB(dbPool, *candleData5Min, address, "5m"); err != nil {
				log.Printf("Ошибка при записи данных 5минутной свечи в файл: %v", err)
			}

			openPrice5Min = 0
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

			if err := WriteCandleToDB(dbPool, *candleData1Hr, address, "1h"); err != nil {
				log.Printf("Ошибка при записи данных часовой свечи в файл: %v", err)
			}

			openPrice1Hr = 0
		}
	}()
}

func WriteCandleToDB(dbPool *pgxpool.Pool, candle CandleData, address, timeframe string) error {
	var tableName string
	switch address {
	case "EQBDMXqg2YcGmMnn5_bXG63y-hh_YNV0dx-ylx-vL3v_WZt4":
		if timeframe == "1h" {
			tableName = "candlesHoursMarketMakers"
		} else if timeframe == "5m" {
			tableName = "candlesMinutesMarketMakers"
		}
	case "EQAl_hUCAeEv-fKtGxYtITAS6PPxuMRaQwHj0QAHeWe6ZSD0":
		if timeframe == "1h" {
			tableName = "candlesHoursLostDogs"
		} else if timeframe == "5m" {
			tableName = "candlesMinutesLostDogs"
		}
	default:
		return fmt.Errorf("неизвестный адрес: %s", address)
	}

	createTableQuery := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			startTime TIMESTAMP NOT NULL,
			endTime TIMESTAMP NOT NULL,
			lowPrice FLOAT NOT NULL,
			highPrice FLOAT NOT NULL,
			open FLOAT NOT NULL,
			close FLOAT NOT NULL
		);`, tableName)

	_, err := dbPool.Exec(context.Background(), createTableQuery)
	if err != nil {
		return fmt.Errorf("ошибка создания таблицы: %v", err)
	}

	insertQuery := fmt.Sprintf(`
		INSERT INTO %s (startTime, endTime, lowPrice, highPrice, open, close)
		VALUES ($1, $2, $3, $4, $5, $6);`, tableName)

	_, err = dbPool.Exec(context.Background(), insertQuery,
		candle.StartTime, candle.EndTime, candle.LowPrice,
		candle.HighPrice, candle.Open, candle.Close)

	if err != nil {
		return fmt.Errorf("ошибка вставки данных: %v", err)
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
