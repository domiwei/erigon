package execmodule

import (
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/erigontech/erigon/execution/builder"
	"github.com/erigontech/erigon/execution/engineapi/engine_helpers"
	"github.com/erigontech/erigon/execution/types"
)

func TestEvictOldBuildersStopsEvicted(t *testing.T) {
	errInterrupted := errors.New("interrupted")
	stopped := make(chan struct{}, engine_helpers.MaxBuilders)

	blockingBuild := func(_ *builder.Parameters, interrupt *atomic.Bool) (*types.BlockWithReceipts, error) {
		for !interrupt.Load() {
			time.Sleep(time.Millisecond)
		}
		stopped <- struct{}{}
		return nil, errInterrupted
	}

	e := &ExecModule{
		builders: make(map[uint64]*builder.BlockBuilder),
	}

	for i := uint64(1); i <= engine_helpers.MaxBuilders; i++ {
		e.builders[i] = builder.NewBlockBuilder(blockingBuild, &builder.Parameters{}, 60)
	}

	e.evictOldBuilders()

	// If evictOldBuilders stopped the evicted builder, its build function
	// returns promptly and signals the channel. Otherwise it keeps running
	// until the 60s max-build-time timer fires.
	select {
	case <-stopped:
	case <-time.After(3 * time.Second):
		t.Fatal("evicted builder was not stopped — goroutine leaked")
	}

	for _, bldr := range e.builders {
		bldr.Stop()
	}
}
