package handlers

import (
	"fmt"
	"log"
	"net/http"
	"prod/internal/nft"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

type WebSocketClient struct {
	conn     *websocket.Conn
	address  string
	sendChan chan nft.CandleData
}

var (
	clients      = make(map[*WebSocketClient]bool)
	clientsMutex = sync.Mutex{}
	upgrader     = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func HandleWebSocketCandleData(w http.ResponseWriter, r *http.Request) {
	collection := chi.URLParam(r, "name")
	log.Printf("Подключение для коллекции: %s", collection)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка при подключении через WebSocket: %v", err), http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	client := &WebSocketClient{
		conn:     conn,
		address:  collection,
		sendChan: make(chan nft.CandleData),
	}

	clientsMutex.Lock()
	clients[client] = true
	clientsMutex.Unlock()

	go func() {
		defer func() {
			clientsMutex.Lock()
			delete(clients, client)
			clientsMutex.Unlock()
		}()
		for candleData := range client.sendChan {
			err := conn.WriteJSON(candleData)
			if err != nil {
				log.Println("Ошибка при отправке данных через WebSocket:", err)
				break
			}
		}
	}()

	for {
		select {
		case candleData := <-getCandleDataChannel(collection):
			log.Printf("Данные успешно получены из канала для коллекции %s: %+v", collection, candleData)

			select {
			case client.sendChan <- candleData:
				log.Printf("Данные отправлены клиенту %s", collection)
			default:
				log.Printf("Канал клиента %s закрыт или переполнен", collection)
				return
			}
		case <-time.After(1 * time.Minute):
			log.Println("Таймаут ожидания новых данных")
			return
		}
	}
}

func getCandleDataChannel(collection string) chan nft.CandleData {
	switch collection {
	case "telegramUsernames":
		return nft.СollectionChannelstelegramUsernamesFloorPriceArray5m
	case "anonymousTelegramNumbers":
		return nft.СollectionChannelsanonymousTelegramNumbersPriceArray5m
	case "tONDNSDomains":
		return nft.СollectionChannelstONDNSDomainsPriceArray5m
	default:
		log.Printf("Неизвестная колzzлекция: %s", collection)
		return nil
	}
}
