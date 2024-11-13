package nft

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type FloorPriceResponse struct {
	Data struct {
		AlphaNftCollectionStats struct {
			FloorPrice float64 `json:"floorPrice"`
		} `json:"alphaNftCollectionStats"`
	} `json:"data"`
}

type FloorPriceData struct {
	Time       string  `json:"time"`
	FloorPrice float64 `json:"floorPrice"`
	Address    string  `json:"address"`
}

var telegramUsernamesFloorPriceArray []FloorPriceData
var anonymousTelegramNumbersPriceArray []FloorPriceData
var tONDNSDomainsPriceArray []FloorPriceData

func GetNFTCollectionFloor(nftCollectionAddress string) (float64, error) {
	query := `query AlphaNftCollectionStats($address: String!) { alphaNftCollectionStats(address: $address) { floorPrice } }`

	reqBody := map[string]interface{}{
		"query": query,
		"variables": map[string]interface{}{
			"address": nftCollectionAddress,
		},
	}

	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return 0, err
	}

	resp, err := http.Post("https://api.getgems.io/graphql", "application/json", bytes.NewBuffer(reqBodyBytes))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("received non-OK HTTP status: %s", resp.Status)
	}

	var floorPriceResp FloorPriceResponse
	if err := json.NewDecoder(resp.Body).Decode(&floorPriceResp); err != nil {
		return 0, err
	}

	return floorPriceResp.Data.AlphaNftCollectionStats.FloorPrice, nil
}

func WriteFloorToArray(floorPrice float64, address string, arrayName string) error {
	newData := FloorPriceData{
		Time:       time.Now().Format("2006-01-02 15:04:05"),
		FloorPrice: floorPrice,
		Address:    address,
	}
	switch arrayName {
	case "telegramUsernamesFloorPriceArray":
		telegramUsernamesFloorPriceArray = append(telegramUsernamesFloorPriceArray, newData)
	case "anonymousTelegramNumbersPriceArray":
		anonymousTelegramNumbersPriceArray = append(anonymousTelegramNumbersPriceArray, newData)
	case "tONDNSDomainsPriceArray":
		tONDNSDomainsPriceArray = append(tONDNSDomainsPriceArray, newData)
	default:
		return fmt.Errorf("неизвестное имя массива: %s", arrayName)
	}

	return nil
}
