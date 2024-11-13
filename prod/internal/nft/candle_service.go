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
	OpenTime  int64   `json:"openTime"`
	CloseTime int64   `json:"closeTime"`
	LowPrice  float64 `json:"lowPrice"`
	HighPrice float64 `json:"highPrice"`
	Open      float64 `json:"open"`
	Close     float64 `json:"close"`
}

var СandleDataTelegramUsernames CandleData
var СandleDataAnonymousTelegramNumbers CandleData
var СandleDataTONDNSDomains CandleData

var Telegram_Usernames_CollectionAddress = "EQCA14o1-VWhS2efqoh_9M1b_A9DtKTuoqfmkn83AbJzwnPi"
var Anonymous_Telegram_Numbers_CollectionAddress = "EQAOQdwdw8kGftJCSFgOErM1mBjYPe4DBPq8-AhF6vr9si5N"
var TON_DNSDomains_CollectionAddress = "EQC3dNlesgVD8YbAazcauIrXBPfiVhMMr5YYk2in0Mtsz0Bz"

var tableMapping = map[string]map[string]string{
	"EQCA14o1-VWhS2efqoh_9M1b_A9DtKTuoqfmkn83AbJzwnPi": {
		"1h": "candlesHoursTelegramUsernamess",
		"5m": "candlesMinutesTelegramUsernames",
	},
	"EQAOQdwdw8kGftJCSFgOErM1mBjYPe4DBPq8-AhF6vr9si5N": {
		"1h": "candlesHoursAnonymousTelegramNumbers",
		"5m": "candlesMinutesAnonymousTelegramNumbers",
	},
	"EQC3dNlesgVD8YbAazcauIrXBPfiVhMMr5YYk2in0Mtsz0Bz": {
		"1h": "candlesHoursTONDNSDomains",
		"5m": "candlesMinutesTONDNSDomains",
	},
}

func getTableName(address, timeframe string) (string, error) {
	if timeframeMap, ok := tableMapping[address]; ok {
		if tableName, ok := timeframeMap[timeframe]; ok {
			return tableName, nil
		}
		return "", fmt.Errorf("неизвестный временной интервал: %s", timeframe)
	}
	return "", fmt.Errorf("неизвестный адрес: %s", address)
}

func UpdateCandleData(address, floorPriceFile, candleFile5Min, candleFile1Hr string, candleData5Min, candleData1Hr *CandleData) {
	var openPrice5Min, openPrice1Hr, closePrice5Min, closePrice1Hr float64
	var openTime5Min, openTime1Hr time.Time

	ctx := context.Background()

	dbPool, err := db.NewClient(ctx, 3, db.DefaultStorageConfig)
	if err != nil {
		log.Fatalf("Ошибка при создании пула соединений: %v", err)
	}

	go func() {
		for {
			floorPrice, err := GetNFTCollectionFloor(address)
			if err != nil {
				log.Printf("Ошибка при получении floor price: %v", err)
			} else {
				if openPrice5Min == 0 {
					openPrice5Min = floorPrice
					openTime5Min = time.Now()
				}
				closePrice5Min = floorPrice

				if openPrice1Hr == 0 {
					openPrice1Hr = floorPrice
					openTime1Hr = time.Now()
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
				OpenTime:  openTime5Min.Unix(),
				CloseTime: time.Now().Unix(),
				LowPrice:  minPrice,
				HighPrice: maxPrice,
				Open:      openPrice5Min,
				Close:     closePrice5Min,
			}

			if err := WriteCandleToDB(dbPool, *candleData5Min, address, "5m"); err != nil {
				log.Printf("Ошибка при записи данных 5минутной свечи в bd: %v", err)
			} else {
				log.Println("Данные 5 мин свечи успешно записаны в bd")
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
				OpenTime:  openTime1Hr.Unix(),
				CloseTime: time.Now().Unix(),
				LowPrice:  minPrice,
				HighPrice: maxPrice,
				Open:      openPrice1Hr,
				Close:     closePrice1Hr,
			}

			if err := WriteCandleToDB(dbPool, *candleData1Hr, address, "1h"); err != nil {
				log.Printf("Ошибка при записи данных часовой свечи в bd: %v", err)
			} else {
				log.Println("Данные часовой свечи успешно записаны в bd")
			}

			openPrice1Hr = 0
		}
	}()
}

func WriteCandleToDB(dbPool *pgxpool.Pool, candle CandleData, address, timeframe string) error {
	tableName, err := getTableName(address, timeframe)
	if err != nil {
		return err
	}

	// Запрос на создание таблицы, если она не существует
	createTableQuery := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s (
		openTime BIGINT NOT NULL,
		closeTime BIGINT NOT NULL,
		lowPrice FLOAT NOT NULL,
		highPrice FLOAT NOT NULL,
		open FLOAT NOT NULL,
		close FLOAT NOT NULL
	);`, tableName)

	_, err = dbPool.Exec(context.Background(), createTableQuery)
	if err != nil {
		return fmt.Errorf("ошибка создания таблицы: %v", err)
	}

	// Запрос на вставку новых данных
	insertQuery := fmt.Sprintf(`
	INSERT INTO %s (openTime, closeTime, lowPrice, highPrice, open, close)
	VALUES ($1, $2, $3, $4, $5, $6);`, tableName)

	_, err = dbPool.Exec(context.Background(), insertQuery,
		candle.OpenTime, candle.CloseTime, candle.LowPrice,
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
