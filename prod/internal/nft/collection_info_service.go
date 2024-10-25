package nft

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Metadata struct {
	SocialLinks []string `json:"social_links"`
	Image       string   `json:"image"`
	Name        string   `json:"name"`
}

type CollectionData struct {
	Metadata Metadata `json:"metadata"`
}

func GetCollectionData(address string) (*CollectionData, error) {
	url := fmt.Sprintf("https://tonapi.io/v2/nfts/collections/%s", address)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var collection CollectionData
	err = json.Unmarshal(body, &collection)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	collection.Metadata.Name = AbbreviateName(collection.Metadata.Name)

	return &collection, nil
}

func AbbreviateName(name string) string {
	words := strings.Fields(name)
	if len(words) <= 1 {
		return name
	}
	var abbreviation string
	for _, word := range words {
		abbreviation += string(word[0])
	}
	return abbreviation
}
