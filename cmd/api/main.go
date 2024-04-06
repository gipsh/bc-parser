package main

import (
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"

	"github.com/gipsh/bc-parser/internal/app"
	"github.com/gipsh/bc-parser/internal/client"
	"github.com/gipsh/bc-parser/internal/indexer"
	"github.com/gipsh/bc-parser/internal/observer"
	"github.com/gipsh/bc-parser/internal/storage"
)

func main() {

	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: false}))

	// eth client
	cli := client.NewEthClient("https://cloudflare-eth.com")

	// storage
	storage := storage.NewStorage()

	// output channel for indexer to send blocks
	output := make(chan client.Block)
	done := make(chan bool, 2)

	var wg sync.WaitGroup
	wg.Add(2)
	// indexer
	indexer := indexer.NewIndexer(log, cli, &wg, done, output, 19590119)

	// observer
	parser := observer.NewObserver(log, &wg, done, output, storage)

	// start indexer and parser
	go indexer.Start()
	go parser.Start()

	// app
	app := app.NewApp(log, parser)

	srv := &http.Server{
		Addr:    ":8081",
		Handler: app.SetupRoutes(),
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("listen", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Info("Shutdown Server ...")

	done <- true
	done <- true
	log.Info("Waiting for indexer and parser to stop")
	wg.Wait()
	log.Info("done. Goodbye!")

}
