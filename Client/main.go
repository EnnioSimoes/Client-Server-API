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

func GetPrice() (string, error) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*100) // 300ms
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "http://127.0.0.1:8080/cotacao", nil)
	if err != nil {
		return "", fmt.Errorf("Erro ao criar requisição: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)

	if ctx.Err() == context.DeadlineExceeded {
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

	var data string

	err = json.Unmarshal(res, &data)

	return data, nil
}

func main() {
	res, err := GetPrice()
	if err != nil {
		log.Println(err)
	}
	println(res)
}
