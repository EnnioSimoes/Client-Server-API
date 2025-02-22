package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type Price struct {
	USDBRL struct {
		Code       string `json:"code"`
		Codein     string `json:"codein"`
		Name       string `json:"name"`
		High       string `json:"high"`
		Low        string `json:"low"`
		VarBid     string `json:"varBid"`
		PctChange  string `json:"pctChange"`
		Bid        string `json:"bid"`
		Ask        string `json:"ask"`
		Timestamp  string `json:"timestamp"`
		CreateDate string `json:"create_date"`
	} `json:"USDBRL"`
}

func GetPrice() (string, error) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*200) // 200ms
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		return "", fmt.Errorf("Erro ao criar requisição: %w", err)
	}
	resp, err := http.DefaultClient.Do(req)

	if ctx.Err() == context.DeadlineExceeded {
		// log.Println("Time to get the price is over")
		return "", ctx.Err()
	}

	if err != nil {
		return "", fmt.Errorf("Erro na requisição HTTP: %w", err)
	}
	defer resp.Body.Close()

	res, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Erro ao ler resposta HTTP: %w", err)
	}

	var data Price

	err = json.Unmarshal(res, &data)
	if err != nil {
		return "", fmt.Errorf("Erro ao decodificar JSON: %w", err)
	}

	return data.USDBRL.Bid, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	price, err := GetPrice()
	if err != nil {
		log.Println(err)
		http.Error(w, "Request cancelled", http.StatusRequestTimeout)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(price)
	// w.Write([]byte(price))
}

func main() {
	http.HandleFunc("/cotacao", handler)
	http.ListenAndServe(":8080", nil)

}
