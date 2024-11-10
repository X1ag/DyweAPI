package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"prod/db"
	"prod/internal/nft"
	"time"

	"github.com/go-chi/chi/v5"
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

	candleDataResponse, err := ReadAllCandlesFromDB(ctx, dbPool, address, timeframe)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка при чтении данных свечей: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(candleDataResponse); err != nil {
		http.Error(w, fmt.Sprintf("Ошибка при кодировании JSON: %v", err), http.StatusInternalServerError)
	}
}

func ReadAllCandlesFromDB(ctx context.Context, dbPool *pgxpool.Pool, address string, timeframe string) ([]nft.CandleData, error) {
	var candles []nft.CandleData
	var query string

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	switch address {
	case "EQBDMXqg2YcGmMnn5_bXG63y-hh_YNV0dx-ylx-vL3v_WZt4":
		if timeframe == "1h" {
			query = "SELECT openTime, closeTime, lowPrice, highPrice, open, close FROM candlesHoursMarketMakers ORDER BY closeTime;"
		} else if timeframe == "5m" {
			query = "SELECT openTime, closeTime, lowPrice, highPrice, open, close FROM candlesMinutesMarketMakers ORDER BY closeTime;"
		}
	case "EQAl_hUCAeEv-fKtGxYtITAS6PPxuMRaQwHj0QAHeWe6ZSD0":
		if timeframe == "1h" {
			query = "SELECT openTime, closeTime, lowPrice, highPrice, open, close FROM candlesHoursLostDogs ORDER BY closeTime;"
		} else if timeframe == "5m" {
			query = "SELECT openTime, closeTime, lowPrice, highPrice, open, close FROM candlesMinutesLostDogs ORDER BY closeTime;"
		}
	default:
		return candles, fmt.Errorf("неизвестный адрес: %s", address)
	}

	rows, err := dbPool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var candle nft.CandleData
		if err := rows.Scan(&candle.OpenTime, &candle.CloseTime, &candle.LowPrice, &candle.HighPrice, &candle.Open, &candle.Close); err != nil {
			return nil, fmt.Errorf("ошибка чтения данных свечи: %v", err)
		}
		candles = append(candles, candle)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("ошибка обработки строк результата: %v", rows.Err())
	}

	return candles, nil
}
