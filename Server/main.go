package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
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

type Cotacao struct {
	Cotacao   string    `json:"cotacao"`
	CreatedAt time.Time `json:"created_at"`
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
		if err == context.DeadlineExceeded {
			http.Error(w, "Request timed out", http.StatusRequestTimeout)
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	cotacao := Cotacao{
		Cotacao:   price,
		CreatedAt: time.Now(),
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*10)
	defer cancel()

	err = Save(cotacao, ctx)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("Timeout ao salvar o produto:", err)
		} else {
			log.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(price)
	// w.Write([]byte(price))
}

func Save(cotacao Cotacao, ctx context.Context) error {
	db, err := sql.Open("sqlite3", "./cotacoes.db")
	if err != nil {
		return fmt.Errorf("Erro ao abrir o banco de dados: %w", err)
	}
	defer db.Close()

	stmt, err := db.PrepareContext(ctx, "INSERT INTO cotacao(cotacao, created_at) VALUES(?, ?)")
	if err != nil {
		return fmt.Errorf("Erro ao preparar a declaração SQL: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, cotacao.Cotacao, cotacao.CreatedAt)
	if err != nil {
		return fmt.Errorf("Erro ao executar a declaração SQL: %w", err)
	}
	return nil
}

func CreateDB(db *sql.DB) error {
	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS cotacao (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            cotacao TEXT NOT NULL,
            created_at TEXT NOT NULL
        );
    `)
	if err != nil {
		return fmt.Errorf("Erro ao criar a tabela: %w", err)
	}

	log.Println("Banco criado com sucesso.")
	return nil
}

func main() {
	db, err := sql.Open("sqlite3", "./cotacoes.db")
	if err != nil {
		log.Println("Erro ao abrir o banco de dados: %w", err)
	}
	defer db.Close()

	err = CreateDB(db)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/cotacao", handler)
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Println("Falha ao subir servidor na porta 8080")
	}
}
