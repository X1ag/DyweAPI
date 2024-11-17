package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

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

	candleData, err := ReadAllCandlesFromDB(ctx, dbPool, client.address, client.timeframe)
	if err != nil {
		log.Println("Ошибка при получении данных свечей:", err)
		conn.WriteJSON(map[string]string{"error": "Failed to fetch candle data"})
		return
	}

	for _, candle := range candleData {
		err := conn.WriteJSON(candle)
		if err != nil {
			log.Println("Ошибка при отправке данных через WebSocket:", err)
			break
		}
	}

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
	lastUpdate := int64(0)

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
			candleData, err := ReadAllCandlesFromDB(ctx, dbPool, client.address, client.timeframe)
			if err != nil {
				log.Println("Ошибка при получении данных свечей:", err)
				return
			}
			var newCandles []nft.CandleData
			for _, candle := range candleData {
				if candle.CloseTime > lastUpdate {
					newCandles = append(newCandles, candle)
				}
			}
			for _, candle := range newCandles {
				log.Printf("Отправка новых данных свечи через WebSocket: %+v\n", candle)

				select {
				case client.sendChan <- candle:
					log.Println("Новые данные свечи успешно отправлены")
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
