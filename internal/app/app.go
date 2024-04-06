package app

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"regexp"

	"github.com/gipsh/bc-parser/internal/client"
	"github.com/gipsh/bc-parser/internal/observer"
)

type App struct {
	log    *slog.Logger
	parser *observer.Observer
	re     *regexp.Regexp
}

func NewApp(log *slog.Logger, parser *observer.Observer) *App {
	return &App{
		log:    log,
		parser: parser,
		re:     regexp.MustCompile("^0x[0-9a-fA-F]{40}$"),
	}
}

type GetLastBlockResponse struct {
	LastBlock int `json:"last_block"`
}

type GetTransactionsResponse struct {
	Transactions []client.Transaction `json:"transactions"`
}

type AddSubscriberResponse struct {
	Status string `json:"status"`
}

// GetLastBlock returns the last block number
func (app *App) GetLastBlock(w http.ResponseWriter, r *http.Request) {
	app.log.Info("App::GetLastBlock")

	lastBlock := app.parser.GetCurrentBlock()

	resp := GetLastBlockResponse{
		LastBlock: lastBlock,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// GetTransactions returns transactions for the given address
func (app *App) GetTransactions(w http.ResponseWriter, r *http.Request) {
	app.log.Info("App::GetTransactions")

	address := r.PathValue("address")

	if !app.isValidAddress(address) {
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

// AddSubscriber adds address to the list of addresses to watch
func (app *App) AddSubscriber(w http.ResponseWriter, r *http.Request) {
	app.log.Info("App::AddSubscriber")

	address := r.PathValue("address")

	if !app.isValidAddress(address) {
		http.Error(w, `{"error":"invalid address"}`, http.StatusBadRequest)
		return
	}

	app.parser.Subscribe(address)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(AddSubscriberResponse{Status: "ok"})

}

func (app *App) SetupRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/v1/last_block", app.GetLastBlock)
	mux.HandleFunc("GET /api/v1/transactions/{address}", app.GetTransactions)
	mux.HandleFunc("POST /api/v1/subscribe/{address}", app.AddSubscriber)

	return mux

}

// IsValidAddress checks if the given string is a valid ethereum address
func (app *App) isValidAddress(v string) bool {
	return app.re.MatchString(v)
}
