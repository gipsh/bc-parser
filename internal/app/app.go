package app

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"regexp"

	"github.com/gipsh/bc-parser/internal/client"
	"github.com/gipsh/bc-parser/internal/observer"
)

type App struct {
	log    *slog.Logger
	parser *observer.Observer
}

func NewApp(log *slog.Logger, parser *observer.Observer) *App {
	return &App{
		log:    log,
		parser: parser,
	}
}

type GetLastBlockResponse struct {
	LastBlock int `json:"last_block"`
}

type GetTransactionsResponse struct {
	Transactions []client.Transaction `json:"transactions"`
}

func (app *App) GetLastBlock(w http.ResponseWriter, r *http.Request) {
	lastBlock := app.parser.GetCurrentBlock()

	resp := GetLastBlockResponse{
		LastBlock: lastBlock,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (app *App) GetTransactions(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")

	if !IsValidAddress(address) {
		http.Error(w, `{"error":"invalid address"}`, http.StatusBadRequest)
		return
	}

	txs := app.parser.GetTransactions(address)

	if txs == nil {
		http.Error(w, `{"error":"address not subscribed"}`, http.StatusBadRequest)
		return
	}

	resp := GetTransactionsResponse{
		Transactions: txs,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (app *App) AddSubscriber(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")

	if !IsValidAddress(address) {
		http.Error(w, `{"error":"invalid address"}`, http.StatusBadRequest)
		return
	}

	app.parser.Subscribe(address)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status": "ok"}`)
}

func (app *App) SetupRoutes() http.Handler {
	mux := http.NewServeMux()

	apiHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {
		case http.MethodGet:
			switch r.URL.Path {
			case "/api/v1/last_block":
				app.GetLastBlock(w, r)
			case "/api/v1/transactions":
				app.GetTransactions(w, r)
			}
		case http.MethodPost:
			app.AddSubscriber(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.Handle("/api/v1/", apiHandler)

	return mux

}

func IsValidAddress(v string) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	return re.MatchString(v)
}
