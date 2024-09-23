package ton

import (
	"log"
	"net/http"

	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/ton" // Убедитесь, что этот импорт добавлен
)

var client *liteclient.ConnectionPool

// InitClient инициализирует соединения с лайтсерверами из global.config.json
func InitClient() {
	client = liteclient.NewConnectionPool()

	// Загружаем конфигурацию
	configUrl := "https://ton-blockchain.github.io/global.config.json"
	resp, err := http.Get(configUrl)
	if err != nil {
		log.Fatalf("Failed to download config: %v", err)
	}
	defer resp.Body.Close()

	// Здесь будет код для подключения к серверам
}

func GetClient() *liteclient.ConnectionPool {
	return client
}

func GetAPIClient() *ton.APIClient {
	return ton.NewAPIClient(client)
}

// package ton

// import (
//     "context"
//     "log"
// 	"fmt"
//     "net/http"
//     "encoding/json"
//     "github.com/xssnick/tonutils-go/liteclient"
//     "golang.org/x/crypto/ed25519"
// )

// var client *liteclient.ConnectionPool

// // Структура для хранения данных из global.config.json
// type Config struct {
//     Liteservers []struct {
//         IP    string `json:"ip"`
//         Port  int    `json:"port"`
//         Key   string `json:"key"`
//     } `json:"liteservers"`
// }

// // InitClient инициализирует соединения с лайтсерверами из global.config.json
// func InitClient() {
//     client = liteclient.NewConnectionPool()

//     // Загружаем конфигурацию
//     configUrl := "https://ton-blockchain.github.io/global.config.json"
//     resp, err := http.Get(configUrl)
//     if err != nil {
//         log.Fatalf("Failed to download config: %v", err)
//     }
//     defer resp.Body.Close()

//     var config Config
//     if err := json.NewDecoder(resp.Body).Decode(&config); err != nil {
//         log.Fatalf("Failed to parse config: %v", err)
//     }

//     // Подключаемся к каждому лайт-серверу
//     for _, server := range config.Liteservers {
//         addr := server.IP + ":" + fmt.Sprint(server.Port)

//         publicKey := server.Key
//         privateKey := ed25519.PrivateKey("your-private-key-here") // Ваш приватный ключ

//         err := client.AddConnection(context.Background(), addr, publicKey, privateKey)
//         if err != nil {
//             log.Fatalf("Failed to connect to %s: %v", addr, err)
//         }
//     }
// }

// func GetClient() *liteclient.ConnectionPool {
//     return client
// }

// ====================================================================================================
// package ton

// import (
// 	"context"
// 	"log"

// 	"github.com/xssnick/tonutils-go/liteclient"
// 	"golang.org/x/crypto/ed25519" //Для работы с ed25519 ключами.
// )

// var client *liteclient.ConnectionPool

// // InitClient инициализирует соединения с лайтсерверами
// func InitClient() {
// 	client = liteclient.NewConnectionPool()
// 	servers := []struct {
// 		address string
// 		port    string
// 	}{
// 		{"testnet.example.com", "8080"}, // !!заглушка!!
// 		{"dev.example.org", "8443"},     // !!заглушка!!
// 	}

// 	for _, server := range servers {
// 		addr := server.address + ":" + server.port

// 		publicKey := "your-public-key-here"
// 		privateKey := ed25519.PrivateKey("your-private-key-here")

// 		err := client.AddConnection(context.Background(), addr, publicKey, privateKey)
// 		if err != nil {
// 			log.Fatalf("Failed to connect to %s: %v", addr, err)
// 		}
// 	}
// }

// func GetClient() *liteclient.ConnectionPool {
// 	return client
// }
