package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"prod/db"
	"prod/internal/nft"

	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v4/pgxpool"
)

type WebSocketClient struct {
	conn      *websocket.Conn
	address   string
	timeframe string
	sendChan  chan nft.CandleData
}

var clients = make(map[*WebSocketClient]bool)
var clientsMutex = sync.Mutex{}
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleWebSocketCandleData(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка при подключении через WebSocket: %v", err), http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	var request struct {
		Collection string `json:"collection"`
		Timeframe  string `json:"timeframe"`
	}
	err = conn.ReadJSON(&request)
	if err != nil {
		log.Println("Ошибка при чтении данных с клиента:", err)
		conn.WriteJSON(map[string]string{"error": "Invalid JSON"})
		return
	}
	validTimeframes := map[string]bool{"5m": true, "1h": true}
	if !validTimeframes[request.Timeframe] {
		conn.WriteJSON(map[string]string{"error": "Invalid timeframe"})
		return
	}

	client := &WebSocketClient{
		conn:      conn,
		address:   request.Collection,
		timeframe: request.Timeframe,
		sendChan:  make(chan nft.CandleData),
	}
	clientsMutex.Lock()
	clients[client] = true
	clientsMutex.Unlock()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbPool, err := db.NewClient(ctx, 3, db.DefaultStorageConfig)
	if err != nil {
		log.Println("Ошибка подключения к базе данных:", err)
		return
	}
	defer dbPool.Close()

	go func() {
		handleCandleUpdates(ctx, client, dbPool)
	}()

	for candleData := range client.sendChan {
		err := conn.WriteJSON(candleData)
		if err != nil {
			log.Println("Ошибка при отправке данных через WebSocket:", err)
			break
		}
	}

	clientsMutex.Lock()
	delete(clients, client)
	clientsMutex.Unlock()
}

func handleCandleUpdates(ctx context.Context, client *WebSocketClient, dbPool *pgxpool.Pool) {
	lastUpdate := time.Now().Unix()

	for {
		select {
		case update := <-nft.UpdateChan:
			parts := strings.Split(update, ":")
			if len(parts) != 2 {
				log.Println("Некорректное уведомление:", update)
				continue
			}
			address, timeframe := parts[0], parts[1]

			if address != client.address || timeframe != client.timeframe {
				continue
			}
			candleData, err := ReadNewCandlesFromDB(ctx, dbPool, client.address, client.timeframe, lastUpdate)
			if err != nil {
				log.Println("Ошибка при получении новых свечей:", err)
				return
			}
			for _, candle := range candleData {
				log.Printf("Отправка новых данных свечи через WebSocket: %+v\n", candle)

				select {
				case client.sendChan <- candle:
					lastUpdate = candle.CloseTime
				case <-ctx.Done():
					log.Println("Контекст отменён, завершение обработки обновлений свечей")
					return
				}
			}
		case <-ctx.Done():
			log.Println("Контекст отменён, завершение обработки обновлений свечей")
			return
		}
	}
}

func ReadNewCandlesFromDB(ctx context.Context, dbPool *pgxpool.Pool, address string, timeframe string, lastUpdate int64) ([]nft.CandleData, error) {
	var candles []nft.CandleData
	var query string

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	switch address {
	case "EQCA14o1-VWhS2efqoh_9M1b_A9DtKTuoqfmkn83AbJzwnPi":
		if timeframe == "1h" {
			query = "SELECT openTime, closeTime, lowPrice, highPrice, open, close FROM candlesHoursTelegramUsernames WHERE closeTime > $1 ORDER BY closeTime;"
		} else if timeframe == "5m" {
			query = "SELECT openTime, closeTime, lowPrice, highPrice, open, close FROM candlesMinutesTelegramUsernames WHERE closeTime > $1 ORDER BY closeTime;"
		}
	case "EQAOQdwdw8kGftJCSFgOErM1mBjYPe4DBPq8-AhF6vr9si5N":
		if timeframe == "1h" {
			query = "SELECT openTime, closeTime, lowPrice, highPrice, open, close FROM candlesHoursAnonymousTelegramNumbers WHERE closeTime > $1 ORDER BY closeTime;"
		} else if timeframe == "5m" {
			query = "SELECT openTime, closeTime, lowPrice, highPrice, open, close FROM candlesMinutesAnonymousTelegramNumbers WHERE closeTime > $1 ORDER BY closeTime;"
		}
	case "EQC3dNlesgVD8YbAazcauIrXBPfiVhMMr5YYk2in0Mtsz0Bz":
		if timeframe == "1h" {
			query = "SELECT openTime, closeTime, lowPrice, highPrice, open, close FROM candlesHoursTONDNSDomains WHERE closeTime > $1 ORDER BY closeTime;"
		} else if timeframe == "5m" {
			query = "SELECT openTime, closeTime, lowPrice, highPrice, open, close FROM candlesMinutesTONDNSDomains WHERE closeTime > $1 ORDER BY closeTime;"
		}
	default:
		return candles, fmt.Errorf("неизвестный адрес: %s", address)
	}
	rows, err := dbPool.Query(ctx, query, lastUpdate)
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
