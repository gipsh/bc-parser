package observer

import (
	"log/slog"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/gipsh/bc-parser/internal/client"
	"github.com/gipsh/bc-parser/internal/common"
	"github.com/gipsh/bc-parser/internal/storage"
)

type Parser interface {
	GetCurrentBlock() int
	Subscribe(address string)
	GetTransactions(address string) []client.Transaction
}

type Observer struct {
	log          *slog.Logger
	input        chan client.Block
	storage      storage.Storer
	currentBlock atomic.Int32
	done         chan bool
	wg           *sync.WaitGroup
}

func NewObserver(log *slog.Logger, wg *sync.WaitGroup, done chan bool, input chan client.Block, storage storage.Storer) *Observer {
	return &Observer{
		log:     log,
		input:   input,
		storage: storage,
		done:    done,
		wg:      wg,
	}
}

func (o *Observer) Subscribe(address string) {
	o.storage.InitAddress(strings.ToLower(address))
}

func (o *Observer) GetTransactions(address string) []client.Transaction {
	return o.storage.GetTx(strings.ToLower(address))
}

func (o *Observer) GetCurrentBlock() int {
	return int(o.currentBlock.Load())
}

func (o *Observer) Start() {
	o.log.Info("Observer started")
	defer o.wg.Done()

	for {
		select {
		case <-o.done:
			o.log.Info("Observer stopped")
			return
		case block := <-o.input:

			o.currentBlock.Store(int32(common.Hex2int(block.Number)))
			o.log.Info("Current block", "block", o.currentBlock.Load(), "tx_count:", len(block.Transactions))

			for _, tx := range block.Transactions {
				// check if transaction is for/to subscriber
				if o.storage.HasAddress(strings.ToLower(tx.From)) {
					o.log.Info("Found From address", "address", tx.From)
					o.storage.Store(tx.From, tx)
				}
				if o.storage.HasAddress(strings.ToLower(tx.To)) {
					o.log.Info("Found To address", "address", tx.To)
					o.storage.Store(tx.To, tx)
				}
			}
		}
	}
}
