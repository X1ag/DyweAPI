package nft

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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

func WriteFloorToFile(floorPrice float64, address string, fileName string) error {
	newData := FloorPriceData{
		Time:       time.Now().Format("2006-01-02 15:04:05"),
		FloorPrice: floorPrice,
		Address:    address,
	}

	var data []FloorPriceData

	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	if fileInfo.Size() > 0 {
		var buf bytes.Buffer
		if _, err := buf.ReadFrom(file); err != nil {
			return err
		}

		if buf.Bytes()[0] == '{' {
			var singleData FloorPriceData
			if err := json.Unmarshal(buf.Bytes(), &singleData); err != nil {
				return err
			}
			data = append(data, singleData)
		} else {
			if err := json.Unmarshal(buf.Bytes(), &data); err != nil {
				return err
			}
		}
	}

	data = append(data, newData)

	file, err = os.OpenFile(fileName, os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		return err
	}

	return nil
}

func GetCandleInfo(fileName string) (float64, float64, error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY, 0644)
	if err != nil {
		return 0, 0, err
	}
	defer file.Close()

	var data []map[string]interface{}
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return 0, 0, err
	}

	if len(data) == 0 {
		return 0, 0, fmt.Errorf("данные отсутствуют")
	}

	minPrice := data[0]["floorPrice"].(float64)
	maxPrice := data[0]["floorPrice"].(float64)

	for _, entry := range data {
		price := entry["floorPrice"].(float64)
		if price < minPrice {
			minPrice = price
		}
		if price > maxPrice {
			maxPrice = price
		}
	}

	return minPrice, maxPrice, nil
}
