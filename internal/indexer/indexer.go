package indexer

import (
	"fmt"
	"log/slog"
	"sync"
	"time"

	c "github.com/gipsh/bc-parser/internal/client"
	"github.com/gipsh/bc-parser/internal/common"
)

type Indexer struct {
	log       *slog.Logger
	client    c.Client
	output    chan c.Block
	done      chan bool
	wg        *sync.WaitGroup
	lastBlock uint64
}

func NewIndexer(log *slog.Logger, client c.Client, wg *sync.WaitGroup, done chan bool, output chan c.Block, startBlock uint64) *Indexer {
	return &Indexer{
		log:       log,
		client:    client,
		output:    output,
		lastBlock: startBlock,
		done:      done,
		wg:        wg,
	}
}

func (i *Indexer) Start() {
	i.log.Info("Indexer started")
	defer i.wg.Done()

	// get last block number
	lb, err := i.client.GetBlockNumber()
	if err != nil {
		i.log.Info("Error getting last block number", err)
		return
	}

	lastBlock := common.Hex2int(lb)
	i.log.Info("Last block number", "lb", lb, "dec", lastBlock)

	// if real last block is less than user start block, set last block to start block
	if i.lastBlock > lastBlock {
		i.lastBlock = lastBlock
	}

	delta := lastBlock - i.lastBlock
	i.log.Info("Last block number", "lb", lastBlock, "Delta", delta)

	for {
		select {
		case <-i.done:
			i.log.Info("Indexer stopped")
			return
		default:
			// if delta is 0, get last block number again and calculate delta
			if delta == 0 {
				delta, err = i.calcDelta()
				if err != nil {
					i.log.Info("Error getting last block number", err)
					return
				}
			}

			// increment last block and decrement delta
			// because only one goroutine is running this is thread safe
			i.lastBlock++
			delta--

			// uint64 to hex string
			blockNumberHex := fmt.Sprintf("%#x", i.lastBlock)
			i.log.Info("Getting block number", "bn", blockNumberHex)

			block, err := i.client.GetBlockByNumber(blockNumberHex)
			if err != nil {
				i.log.Error("Error getting block number", "error", err)
				continue
			}

			// send block to output channel
			i.output <- *block

			// sleep 5scs to avoid rate limiting
			time.Sleep(5 * time.Second)

		}
	}
}

// calcDelta calculates the difference between the last block number and
// the last block number stored in the indexer (init block)
func (i *Indexer) calcDelta() (uint64, error) {

	lb, err := i.client.GetBlockNumber()
	if err != nil {
		i.log.Info("Error getting last block number", err)
		return 0, err
	}

	lastBlock := common.Hex2int(lb)
	delta := lastBlock - i.lastBlock
	return delta, nil
}
