package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"prod/db"
	"prod/internal/nft"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v4/pgxpool"
)

func HandleCandleData(w http.ResponseWriter, r *http.Request) {
	address := chi.URLParam(r, "address")
	timeframe := chi.URLParam(r, "timeframe")

	ctx := context.Background()

	dbPool, err := db.NewClient(ctx, 3, db.StorageConfig{})
	if err != nil {
		log.Fatalf("Ошибка при создании пула соединений: %v", err)
	}
	defer dbPool.Close()

	var candleDataResponse nft.CandleData

	candleDataResponse, err = ReadLastCandleFromDB(ctx, dbPool, address, timeframe)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading candle data: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(candleDataResponse)
}

func ReadLastCandleFromDB(ctx context.Context, dbPool *pgxpool.Pool, address string, timeframe string) (nft.CandleData, error) {
	var candle nft.CandleData

	// Определение SQL-запроса в зависимости от адреса и временного интервала
	var query string
	switch address {
	case "EQBDMXqg2YcGmMnn5_bXG63y-hh_YNV0dx-ylx-vL3v_WZt4":
		if timeframe == "1h" {
			query = "SELECT startTime, endTime, lowPrice, highPrice, open, close FROM candlesHoursMarketMakers ORDER BY endTime DESC LIMIT 1;"
		} else if timeframe == "5m" {
			query = "SELECT startTime, endTime, lowPrice, highPrice, open, close FROM candlesMinutesMarketMakers ORDER BY endTime DESC LIMIT 1;"
		}
	case "EQAl_hUCAeEv-fKtGxYtITAS6PPxuMRaQwHj0QAHeWe6ZSD0":
		if timeframe == "1h" {
			query = "SELECT startTime, endTime, lowPrice, highPrice, open, close FROM candlesHoursLostDogs ORDER BY endTime DESC LIMIT 1;"
		} else if timeframe == "5m" {
			query = "SELECT startTime, endTime, lowPrice, highPrice, open, close FROM candlesMinutesLostDogs ORDER BY endTime DESC LIMIT 1;"
		}
	default:
		return candle, fmt.Errorf("неизвестный адрес: %s", address)
	}

	// Выполнение запроса
	row := dbPool.QueryRow(ctx, query)

	// Считывание данных в структуру candle
	err := row.Scan(&candle.StartTime, &candle.EndTime, &candle.LowPrice, &candle.HighPrice, &candle.Open, &candle.Close)
	if err != nil {
		return candle, fmt.Errorf("ошибка считывания данных: %v", err)
	}

	return candle, nil
}
