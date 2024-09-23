// package services

// import (
// 	"context"
// 	"fmt"
// 	"log"

// 	"dywego/internal/ton"
// 	"dywego/pkg/models"
// 	"github.com/xssnick/tonutils-go/liteclient"
// 	"github.com/xssnick/tonutils-go/nft"
// )

package services

import (
	"context"
	"dywego/internal/ton"
	"fmt"

	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/ton/nft"
)

// GetCollection получает название коллекции по адресу
func GetCollection(collectionAddress string) (string, error) {
	// Преобразуем строковый адрес в адрес TON
	addr := address.MustParseAddr(collectionAddress)

	// Получаем API клиент TON
	api := ton.GetAPIClient()

	// Создаём клиент коллекции
	collectionClient := nft.NewCollectionClient(api, addr)

	// Получаем данные о коллекции
	collectionData, err := collectionClient.GetCollectionData(context.Background())
	if err != nil {
		return "", fmt.Errorf("failed to get collection data: %v", err)
	}

	// Получаем название коллекции
	var collectionName string
	switch content := collectionData.Content.(type) {
	case *nft.ContentOffchain:
		collectionName = content.URI // или любой другой атрибут, содержащий название
	case *nft.ContentOnchain:
		collectionName = content.GetAttribute("name")
	default:
		collectionName = "Unknown Collection"
	}

	return collectionName, nil
}

// package services

// import (
// 	"context"
// 	"fmt"
// 	"dywego/internal/ton"
// 	"github.com/xssnick/tonutils-go/address"
// 	"github.com/xssnick/tonutils-go/ton/nft"
// )

// // GetCollection получает название коллекции по адресу
// func GetCollection(collectionAddress string) (string, error) {
// 	// Преобразуем строковый адрес в адрес TON
// 	addr := address.MustParseAddr(collectionAddress)

// 	// Получаем клиент TON
// 	api := ton.GetClient()

// 	// Создаём клиент коллекции
// 	collectionClient := nft.NewCollectionClient(api, addr)

// 	// Получаем данные о коллекции
// 	collectionData, err := collectionClient.GetCollectionData(context.Background())
// 	if err != nil {
// 		return "", fmt.Errorf("failed to get collection data: %v", err)
// 	}

// 	// Получаем название коллекции в зависимости от типа контента
// 	var collectionName string
// 	switch content := collectionData.Content.(type) {
// 	case *nft.ContentOffchain:
// 		collectionName = content.URI // или любой другой атрибут, содержащий название
// 	case *nft.ContentOnchain:
// 		collectionName = content.GetAttribute("name")
// 	default:
// 		collectionName = "Unknown Collection"
// 	}

// 	return collectionName, nil
// }

// =========================================================================================
// type CollectionService struct {
// 	TonClient *liteclient.ConnectionPool // Используем ConnectionPool
// }

// // GetCollection получает информацию о коллекции по адресу контракта
// func (s *CollectionService) GetCollection(ctx context.Context, collectionAddress string) (*models.Collection, error) {
// 	client := ton.NewAPIClient(s.TonClient)

// 	// Создаем клиент для работы с коллекцией
// 	collectionClient := nft.NewCollectionClient(client, collectionAddress)

// 	// Получаем данные о коллекции
// 	collectionData, err := collectionClient.GetCollectionData(ctx)
// 	if err != nil {
// 		log.Printf("Failed to get collection data: %v", err)
// 		return nil, err
// 	}

// 	// Извлекаем имя и другие атрибуты коллекции
// 	collectionName := collectionData.Content.(*nft.ContentOffchain).URI // Предполагаем, что это имя
// 	collectionOwner := collectionData.OwnerAddress.String()
// 	mintedItemsNum := collectionData.NextItemIndex

// 	// Печатаем информацию о коллекции (по желанию)
// 	fmt.Println("Collection addr      :", collectionAddress)
// 	fmt.Println("    content          :", collectionName)
// 	fmt.Println("    owner            :", collectionOwner)
// 	fmt.Println("    minted items num :", mintedItemsNum)

// 	// Возвращаем структуру Collection
// 	return &models.Collection{
// 		ID:   collectionAddress,
// 		Name: collectionName,
// 	}, nil
// }
// ================================================================================================

// package services

// import (
//     "context"
//     "fmt"
//     "dywego/internal/ton"
//     "github.com/xssnick/tonutils-go/liteclient"
//     "github.com/xssnick/tonutils-go/tlb" // Для работы с данными TL-блоков
//     "github.com/xssnick/tonutils-go/tvm/cell" // Для работы с ячейками
//     "dywego/pkg/models"
// )

// type CollectionService struct {
//     TonClient *liteclient.ConnectionPool // Используем ConnectionPool
// }

// // GetCollection получает только имя коллекции по ID
// func (s *CollectionService) GetCollection(ctx context.Context, collectionID string) (*models.Collection, error) {
//     // Для начала получаем данные о коллекции с помощью контракта
//     contractAddress := "адрес контракта коллекции NFT" // Здесь должен быть реальный адрес контракта коллекции

//     // Запросить данные через клиента
//     master, err := ton.GetClient().GetMasterchainInfo(ctx)
//     if err != nil {
//         return nil, fmt.Errorf("failed to get masterchain info: %w", err)
//     }

//     block, err := ton.GetClient().GetBlockData(ctx, master.LastBlockID)
//     if err != nil {
//         return nil, fmt.Errorf("failed to get block data: %w", err)
//     }

//     // Создаём ячейку с идентификатором коллекции
//     collectionCell := cell.BeginCell().MustStoreString(collectionID).EndCell()

//     // Отправляем запрос в блокчейн для получения информации о коллекции
//     result, err := ton.GetClient().RunGetMethod(ctx, block, contractAddress, "get_collection_data", collectionCell)
//     if err != nil {
//         return nil, fmt.Errorf("failed to get collection data: %w", err)
//     }

//     // Парсим имя коллекции из возвращённых данных
//     collectionName, err := result.Int(0) // Допустим, что имя коллекции — это первый возвращённый элемент
//     if err != nil {
//         return nil, fmt.Errorf("failed to parse collection name: %w", err)
//     }

//     // Возвращаем структуру с ID и именем коллекции
//     return &models.Collection{
//         ID:   collectionID,
//         Name: collectionName.String(), // Конвертируем имя коллекции в строку
//     }, nil
// }
