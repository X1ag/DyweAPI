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
	conn      *websocket.Conn
	name      string
	timeframe string // Добавляем поле для временного интервала
	sendChan  chan nft.CandleData
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
	timeframe := chi.URLParam(r, "timeframe")
	log.Printf("Подключение для коллекции: %s с временным интервалом: %s", collection, timeframe)

	// Устанавливаем соединение WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка при подключении через WebSocket: %v", err), http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	// Создаем нового клиента с учетом временного интервала
	client := &WebSocketClient{
		conn:      conn,
		name:      collection,
		timeframe: timeframe, // Сохраняем временной интервал
		sendChan:  make(chan nft.CandleData),
	}

	// Добавляем клиента в список
	clientsMutex.Lock()
	clients[client] = true
	clientsMutex.Unlock()

	// Горутинка для отправки данных клиенту
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

	// Цикл получения данных с канала
	for {
		select {
		case candleData := <-getCandleDataChannel(collection, timeframe):
			log.Printf("Данные успешно получены из канала для коллекции %s и интервала %s: %+v", collection, timeframe, candleData)

			// Отправка данных клиенту через WebSocket
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

func getCandleDataChannel(collection, timeframe string) chan nft.CandleData {
	// Выбираем канал в зависимости от временного интервала
	switch collection {
	case "telegramUsernames":
		switch timeframe {
		case "1h":
			return nft.СollectionChannelstelegramUsernamesFloorPriceArray1h
		case "5m":
			return nft.СollectionChannelstelegramUsernamesFloorPriceArray5m
		case "15m":
			return nft.СollectionChannelstelegramUsernamesFloorPriceArray15m
		case "30m":
			return nft.СollectionChannelstelegramUsernamesFloorPriceArray30m
		case "4h":
			return nft.СollectionChannelstelegramUsernamesFloorPriceArray4h
		default:
			log.Printf("Неизвестный временной интервал %s для коллекции %s", timeframe, collection)
			return nil
		}
	case "anonymousTelegramNumbers":
		switch timeframe {
		case "1h":
			return nft.СollectionChannelsanonymousTelegramNumbersPriceArray1h
		case "5m":
			return nft.СollectionChannelsanonymousTelegramNumbersPriceArray5m
		case "15m":
			return nft.СollectionChannelsanonymousTelegramNumbersPriceArray15m
		case "30m":
			return nft.СollectionChannelsanonymousTelegramNumbersPriceArray30m
		case "4h":
			return nft.СollectionChannelsanonymousTelegramNumbersPriceArray4h
		default:
			log.Printf("Неизвестный временной интервал %s для коллекции %s", timeframe, collection)
			return nil
		}
	case "tONDNSDomains":
		switch timeframe {
		case "1h":
			return nft.СollectionChannelstONDNSDomainsPriceArray1h
		case "5m":
			return nft.СollectionChannelstONDNSDomainsPriceArray5m
		case "15m":
			return nft.СollectionChannelstONDNSDomainsPriceArray15m
		case "30m":
			return nft.СollectionChannelstONDNSDomainsPriceArray30m
		case "4h":
			return nft.СollectionChannelstONDNSDomainsPriceArray4h
		default:
			log.Printf("Неизвестный временной интервал %s для коллекции %s", timeframe, collection)
			return nil
		}
	default:
		log.Printf("Неизвестная коллекция: %s", collection)
		return nil
	}
}
