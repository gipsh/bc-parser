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

	parser.Subscribe("0x1c9fce6dd765a22040d500019ada91acce65b5d2")
	parser.Subscribe("0x6907894f656b95d67e380349a5edc1f75bc45b8c")

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

	quit := make(chan os.Signal, 2)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Info("Shutdown Server ...")

	// stop indexer and parser
	done <- true
	done <- true
	log.Info("Waiting for indexer and parser to stop")
	wg.Wait()
	log.Info("done. Goodbye!")

}
