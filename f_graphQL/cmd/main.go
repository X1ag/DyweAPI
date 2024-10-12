package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

type FloorPriceResponse struct {
	Data struct {
		AlphaNftCollectionStats struct {
			FloorPrice float64 `json:"floorPrice"`
		} `json:"alphaNftCollectionStats"`
	} `json:"data"`
}

type Metadata struct {
	Name string `json:"name"`
}

type CollectionResponse struct {
	Metadata Metadata `json:"metadata"`
}

func GetNFTCollectionFloor(nftCollectionAddress string) (float64, error) {
	query := `
		query AlphaNftCollectionStats($address: String!) {
			alphaNftCollectionStats(address: $address) {
				floorPrice
			}
		}
	`

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

var collectionType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Collection",
		Fields: graphql.Fields{
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"floorPrice": &graphql.Field{
				Type: graphql.Float,
			},
		},
	},
)

var rootQuery = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"collection": &graphql.Field{
				Type: collectionType,
				Args: graphql.FieldConfigArgument{
					"address": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					address := params.Args["address"].(string)

					client := resty.New()

					resp, err := client.R().
						Get(fmt.Sprintf("https://tonapi.io/v2/nfts/collections/%s", address))
					if err != nil || resp.StatusCode() != http.StatusOK {
						return nil, fmt.Errorf("Ошибка запроса к TONAPI")
					}

					var collection CollectionResponse
					if err := json.Unmarshal(resp.Body(), &collection); err != nil {
						return nil, fmt.Errorf("Ошибка парсинга ответа")
					}

					floorPrice, err := GetNFTCollectionFloor(address)
					if err != nil {
						return nil, fmt.Errorf("Ошибка при получении floor price: %v", err)
					}

					return map[string]interface{}{
						"name":       collection.Metadata.Name,
						"floorPrice": floorPrice,
					}, nil
				},
			},
		},
	},
)

func main() {
	schema, err := graphql.NewSchema(
		graphql.SchemaConfig{
			Query: rootQuery,
		},
	)
	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}

	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true, 
	})
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.Handle("/dywe/query", h)

	fmt.Printf("Сервер запущен на порте %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// query {
// 	collection(address: "EQBFg46ihgN95_3Ld7MU19kVJdKepJ0Dq3UHRBaEuJLlomQI") {
// 		name
// 		floorPrice
// 	  }
// 	}
