package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/gorilla/mux"
)

type Metadata struct {
	Name string `json:"name"`
}

type CollectionResponse struct {
	Metadata Metadata `json:"metadata"`
}

func getCollectionName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]

	client := resty.New()

	resp, err := client.R().
		Get(fmt.Sprintf("https://tonapi.io/v2/nfts/collections/%s", address))
	if err != nil {
		http.Error(w, "Ошибка запроса к TONAPI", http.StatusInternalServerError)
		return
	}

	log.Println("Тело ответа от TONAPI:", resp.String())

	if resp.StatusCode() != http.StatusOK {
		http.Error(w, "Ошибка при запросе к TONAPI", http.StatusInternalServerError)
		return
	}

	var collection CollectionResponse
	if err := json.Unmarshal(resp.Body(), &collection); err != nil {
		http.Error(w, "Ошибка парсинга ответа", http.StatusInternalServerError)
		return
	}

	if collection.Metadata.Name == "" {
		http.Error(w, "Коллекция не найдена", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"name": collection.Metadata.Name,
	})
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/dywe/api/v1/collection/{address}", getCollectionName).Methods("GET")

	fmt.Println("Сервер запущен на порте 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
