package indexer

import (
	"log/slog"
	"os"
	"sync"
	"testing"

	"github.com/gipsh/bc-parser/internal/client"
	m "github.com/gipsh/bc-parser/internal/client/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestIndexer_Delta(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cli := m.NewMockClient(ctrl)

	cli.EXPECT().GetBlockNumber().Return("0x12ae9d6", nil).Times(1)

	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: false}))

	wg := sync.WaitGroup{}
	wg.Add(1)
	done := make(chan bool, 1)
	output := make(chan client.Block)
	i := NewIndexer(log, cli, &wg, done, output, 19589580)

	x, err := i.calcDelta()

	assert.Nil(t, err)
	assert.Equal(t, uint64(10), x)

}
