package nft

import (
	"context"
	"fmt"
	"log"
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
		"1h": "candlesHoursTelegramUsernames",
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

func UpdateCandleData(address, floorPriceArray5m, floorPriceArray1h string, candleData5Min, candleData1Hr *CandleData) {
	var openPrice5Min, openPrice1Hr, closePrice5Min, closePrice1Hr float64
	var openTime5Min, openTime1Hr time.Time

	ctx := context.Background()
	dbPool, err := db.NewClient(ctx, 3, db.DefaultStorageConfig)
	if err != nil {
		log.Fatalf("Ошибка при создании пула соединений: %v", err)
	}

	go func() {
		WaitUntilNextInterval(5 * time.Minute)
		for {
			floorPrice, err := GetNFTCollectionFloor(address)
			if err != nil {
				log.Printf("Ошибка при получении floor price: %v", err)
			} else {
				log.Printf("Получен floor price для коллекции %s: %f", address, floorPrice)

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

				if err := WriteFloorToArray(floorPrice, address, floorPriceArray5m); err != nil {
					log.Printf("Ошибка при записи в 5-минутный массив: %v", err)
				}

				if err := WriteFloorToArray(floorPrice, address, floorPriceArray1h); err != nil {
					log.Printf("Ошибка при записи в 1-часовой массив: %v", err)
				}
			}
			time.Sleep(20 * time.Second)
		}
	}()

	go func() {
		WaitUntilNextInterval(5 * time.Minute)
		for {
			time.Sleep(5 * time.Minute)
			minPrice, maxPrice, err := GetCandleInfoFromArray(floorPriceArray5m)
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
				log.Printf("Ошибка при записи данных 5-минутной свечи в базу данных: %v", err)
			} else {
				log.Printf("Данные 5-минутной свечи успешно записаны в базу данных: %+v", *candleData5Min)
			}

			openPrice5Min = 0
			closePrice5Min = 0
			ClearArray(floorPriceArray5m)
		}
	}()

	go func() {
		WaitUntilNextInterval(time.Hour)
		for {
			time.Sleep(1 * time.Hour)
			minPrice, maxPrice, err := GetCandleInfoFromArray(floorPriceArray1h)
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
				log.Printf("Ошибка при записи данных часовой свечи в базу данных: %v", err)
			} else {
				log.Printf("Данные часовой свечи успешно записаны в базу данных: %+v", *candleData1Hr)
			}

			openPrice1Hr = 0
			closePrice1Hr = 0
			ClearArray(floorPriceArray1h)
		}
	}()
}

func WriteCandleToDB(dbPool *pgxpool.Pool, candle CandleData, address, timeframe string) error {
	tableName, err := getTableName(address, timeframe)
	if err != nil {
		return err
	}
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

func GetCandleInfoFromArray(arrayName string) (float64, float64, error) {
	arrays := map[string][]FloorPriceData{
		"telegramUsernamesFloorPriceArray5m":   telegramUsernamesFloorPriceArray5m,
		"anonymousTelegramNumbersPriceArray5m": anonymousTelegramNumbersPriceArray5m,
		"tONDNSDomainsPriceArray5m":            tONDNSDomainsPriceArray5m,
		"telegramUsernamesFloorPriceArray1h":   telegramUsernamesFloorPriceArray1h,
		"anonymousTelegramNumbersPriceArray1h": anonymousTelegramNumbersPriceArray1h,
		"tONDNSDomainsPriceArray1h":            tONDNSDomainsPriceArray1h,
	}

	floorPriceArray, exists := arrays[arrayName]
	if !exists {
		return 0, 0, fmt.Errorf("неизвестное имя массива: %s", arrayName)
	}

	if len(floorPriceArray) == 0 {
		return 0, 0, fmt.Errorf("данные отсутствуют")
	}

	minPrice := floorPriceArray[0].FloorPrice
	maxPrice := floorPriceArray[0].FloorPrice

	for _, entry := range floorPriceArray {
		price := entry.FloorPrice
		if price < minPrice {
			minPrice = price
		}
		if price > maxPrice {
			maxPrice = price
		}
	}

	return minPrice, maxPrice, nil
}

func WaitUntilNextInterval(interval time.Duration) {
	now := time.Now()
	var next time.Time

	if interval == time.Hour {
		next = now.Truncate(time.Hour).Add(time.Hour)
	} else if interval == 5*time.Minute {
		next = now.Truncate(5 * time.Minute).Add(5 * time.Minute)
	} else {
		next = now.Truncate(interval).Add(interval)
	}

	time.Sleep(time.Until(next))
}

func ClearArray(arrayName string) {
	priceArrays := map[string]*[]FloorPriceData{
		"telegramUsernamesFloorPriceArray1h":   &telegramUsernamesFloorPriceArray1h,
		"anonymousTelegramNumbersPriceArray1h": &anonymousTelegramNumbersPriceArray1h,
		"tONDNSDomainsPriceArray1h":            &tONDNSDomainsPriceArray1h,
		"telegramUsernamesFloorPriceArray5m":   &telegramUsernamesFloorPriceArray5m,
		"anonymousTelegramNumbersPriceArray5m": &anonymousTelegramNumbersPriceArray5m,
		"tONDNSDomainsPriceArray5m":            &tONDNSDomainsPriceArray5m,
	}
	if arr, ok := priceArrays[arrayName]; ok {
		*arr = nil
		log.Printf("Массив %s очищен", arrayName)
	} else {
		log.Printf("Не найден массив с именем %s", arrayName)
	}
}
