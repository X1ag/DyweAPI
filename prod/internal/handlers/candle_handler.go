package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"prod/db"
	"prod/internal/nft"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v4/pgxpool"
)

func HandleCandleData(w http.ResponseWriter, r *http.Request) {
	address := chi.URLParam(r, "address")
	timeframe := chi.URLParam(r, "timeframe")

	ctx := context.Background()

	dbPool, err := db.NewClient(ctx, 3, db.DefaultStorageConfig)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка при создании пула соединений: %v", err), http.StatusInternalServerError)
		return
	}
	if dbPool == nil {
		http.Error(w, "Не удалось инициализировать пул соединений", http.StatusInternalServerError)
		return
	}

	candleDataResponse, err := ReadLastCandleFromDB(ctx, dbPool, address, timeframe)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка при чтении данных свечи: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(candleDataResponse); err != nil {
		http.Error(w, fmt.Sprintf("Ошибка при кодировании JSON: %v", err), http.StatusInternalServerError)
	}
}

func ReadLastCandleFromDB(ctx context.Context, dbPool *pgxpool.Pool, address string, timeframe string) (nft.CandleData, error) {
	var candle nft.CandleData
	var query string

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

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

	row := dbPool.QueryRow(ctx, query)

	err := row.Scan(&candle.StartTime, &candle.EndTime, &candle.LowPrice, &candle.HighPrice, &candle.Open, &candle.Close)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Printf("Запись не найдена для адреса %s и временного интервала %s", address, timeframe)
			return candle, nil 
		}

	}

	return candle, nil
}
